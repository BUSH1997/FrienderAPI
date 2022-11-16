package http

import (
	"encoding/json"
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/api/errors/convert"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/chat"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type ChatHandler struct {
	useCase   chat.Usecase
	messenger *chat.Messenger
	logger    hardlogger.Logger
}

func NewChatHandler(usecase chat.Usecase, messenger *chat.Messenger, logger hardlogger.Logger) *ChatHandler {
	return &ChatHandler{
		useCase:   usecase,
		messenger: messenger,
		logger:    logger,
	}
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (ch *ChatHandler) GetChats(echoCtx echo.Context) error {
	ctx := ch.logger.WithCaller(echoCtx.Request().Context())

	chats, err := ch.useCase.GetChats(ctx)
	if err != nil {
		ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to get chats")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, chats)
}

func (ch *ChatHandler) GetMessages(echoCtx echo.Context) error {
	ctx := ch.logger.WithCaller(echoCtx.Request().Context())

	opts := models.GetMessageOpts{}

	pageString := echoCtx.QueryParam("page")
	if pageString != "" {
		page, err := strconv.ParseInt(pageString, 10, 32)
		if err != nil {
			ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to parse page param")
			return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
		}

		opts.Page = int(page)
	}

	limitString := echoCtx.QueryParam("limit")
	if limitString != "" {
		limit, err := strconv.ParseInt(limitString, 10, 32)
		if err != nil {
			ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to parse limit param")
			return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
		}

		opts.Limit = int(limit)
	}

	opts.EventID = echoCtx.QueryParam("event")

	messages, err := ch.useCase.GetMessages(ctx, opts)
	if err != nil {
		ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to get messages")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	user := context.GetUser(ctx)
	err = ch.useCase.UpdateLastCheckTime(ctx, opts.EventID, user, time.Now().Unix())
	if err != nil {
		ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to update last check time for %d", user)
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, messages)
}

func (ch *ChatHandler) ProcessMessage(echoCtx echo.Context) error {
	ctx := ch.logger.WithCaller(echoCtx.Request().Context())

	eventID := echoCtx.Param("id")
	if eventID == "" {
		err := errors.New("event id is empty")
		ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to get event id")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	ws, err := upgrader.Upgrade(echoCtx.Response(), echoCtx.Request(), nil)
	if err != nil {
		ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to upgrade http request")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	user := context.GetUser(ctx)

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
			ch.logger.WithCtx(ctx).WithError(err).Error("failed to read message from socket")
			return echoCtx.JSON(http.StatusInternalServerError, errors.Wrap(err, "failed to read message").Error())
		}

		message := models.Message{
			UserID:      user,
			EventID:     eventID,
			Text:        string(msg),
			TimeCreated: time.Now().Unix(),
		}

		err = ch.useCase.CreateMessage(ctx, message)
		if err != nil {
			ch.logger.WithCtx(ctx).WithError(err).Error("failed to create message")
			return echoCtx.JSON(http.StatusInternalServerError, errors.Wrap(err, "failed to create message").Error())
		}

		jsonMessage, err := json.Marshal(&message)
		if err != nil {
			return echoCtx.JSON(http.StatusInternalServerError, errors.Wrap(err, "failed to unmarshal json message").Error())
		}

		for _, client := range ch.messenger.Chat.Clients {
			fmt.Printf("write to client %d message: %s \n", client.UserID, msg)

			err := client.Socket.WriteMessage(websocket.TextMessage, jsonMessage)
			if err != nil {
				ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to write message to %d", client.UserID)
				ch.messenger.Chat.Clients = remove(ch.messenger.Chat.Clients, user)
				continue
			}

			err = ch.useCase.UpdateLastCheckTime(ctx, eventID, client.UserID, time.Now().Unix())
			if err != nil {
				ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to update last check time for %d", client.UserID)
				return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
			}
		}
	}
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
