package http

import (
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/chat"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type ChatHandler struct {
	useCase   chat.Usecase
	messenger *chat.Messenger
	logger    *logrus.Logger
}

func NewChatHandler(usecase chat.Usecase, messenger *chat.Messenger, logger *logrus.Logger) *ChatHandler {
	return &ChatHandler{
		useCase:   usecase,
		messenger: messenger,
		logger:    logger,
	}
}

var (
	upgrader = websocket.Upgrader{}
)

func (ch *ChatHandler) GetChats(ctx echo.Context) error {
	chats, err := ch.useCase.GetChats(ctx.Request().Context())
	if err != nil {
		ch.logger.WithError(err).Errorf("failed to get chats")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, chats)
}

func (ch *ChatHandler) GetMessages(ctx echo.Context) error {
	opts := models.GetMessageOpts{}

	pageString := ctx.QueryParam("page")
	if pageString != "" {
		page, err := strconv.ParseInt(pageString, 10, 32)
		if err != nil {
			ch.logger.WithError(err).Errorf("failed to parse page param")
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}

		opts.Page = int(page)
	}

	limitString := ctx.QueryParam("limit")
	if limitString != "" {
		limit, err := strconv.ParseInt(limitString, 10, 32)
		if err != nil {
			ch.logger.WithError(err).Errorf("failed to parse limit param")
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}

		opts.Limit = int(limit)
	}

	opts.EventID = ctx.QueryParam("event")

	messages, err := ch.useCase.GetMessages(ctx.Request().Context(), opts)
	if err != nil {
		ch.logger.WithError(err).Errorf("failed to get messages")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, messages)
}

func (ch *ChatHandler) ProcessMessage(ctx echo.Context) error {
	eventID := ctx.Param("id")
	if eventID == "" {
		err := errors.New("event id is empty")
		ch.logger.WithError(err).Errorf("failed to get event id")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		ch.logger.WithError(err).Errorf("failed to upgrade http request")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	user := context.GetUser(ctx.Request().Context())

	ch.messenger.Chat.Mx.Lock()

	ch.messenger.Chat.Clients = append(ch.messenger.Chat.Clients, &chat.Client{
		Socket: ws,
		UserID: user,
	})

	fmt.Printf("client %d connected\n", user)

	ch.messenger.Chat.Mx.Unlock()

	defer func() {
		ws.Close()
		ch.messenger.Chat.Mx.Lock()
		ch.messenger.Chat.Clients = remove(ch.messenger.Chat.Clients, user)
		ch.messenger.Chat.Mx.Unlock()
	}()

	for {
		mt, msg, err := ws.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			ch.logger.WithError(err).Error("failed to read message from socket")
			break
			// TODO: return
		}

		err = ch.useCase.CreateMessage(ctx.Request().Context(), models.Message{
			UserID:      user,
			EventID:     eventID,
			Text:        string(msg),
			TimeCreated: time.Now().Unix(),
		})
		if err != nil {
			ch.logger.WithError(err).Error("failed to create message")
			return ctx.JSON(http.StatusInternalServerError, errors.Wrap(err, "failed to create message").Error())
		}

		for _, client := range ch.messenger.Chat.Clients {
			fmt.Printf("write to client %d message: %s \n", client.UserID, msg)

			err := client.Socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				ch.logger.WithError(err).Errorf("failed to write message to %d", client.UserID)
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
		}
	}

	return ctx.JSON(http.StatusOK, nil)
}

func remove(clients []*chat.Client, uid int64) []*chat.Client {
	for i, client := range clients {
		if client.UserID != uid {
			continue
		}

		return append(clients[:i], clients[i+1:]...)
	}

	return nil
}
