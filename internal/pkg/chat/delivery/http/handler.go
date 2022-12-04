package http

import (
	"context"
	"encoding/json"
	"github.com/BUSH1997/FrienderAPI/internal/api/errors/convert"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/chat"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/gofrs/uuid"
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

	user := contextlib.GetUser(ctx)
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

	user := contextlib.GetUser(ctx)

	if !ch.messenger.HasChat(eventID) {
		ch.messenger.AppendChat(eventID)
	}

	ch.messenger.AppendClientToChat(eventID, ws, user)

	defer func() {
		err := ws.Close()
		if err != nil {
			ch.logger.WithCtx(ctx).WithError(err).Error("failed to close message from socket")
			return // TODO:
		}

		ch.messenger.RemoveClientFromChat(eventID, user)
	}()

	for {
		msg := models.MessageInput{}
		err = ws.ReadJSON(&msg)
		if err != nil {
			ch.logger.WithCtx(ctx).WithError(err).Error("failed to read json message from websocket")
			return echoCtx.JSON(http.StatusInternalServerError, errors.Wrap(err, "failed to read json message from websocket").Error())
		}

		err = ch.routeProcessMessage(ctx, msg, user, eventID)
		if err != nil {
			ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to process message, type %s", msg.Type)
			return echoCtx.JSON(http.StatusInternalServerError, errors.Wrapf(err, "failed to process message, , type %s", msg.Type).Error())
		}
	}
}

func (ch *ChatHandler) routeProcessMessage(ctx context.Context, msg models.MessageInput, user int64, event string) error {
	if msg.Type == chat.CreateTextMessage {
		return ch.sendTextMessage(ctx, msg, user, event)
	}
	if msg.Type == chat.DeleteMessage {
		return ch.deleteMessage(ctx, msg, user, event)
	}

	return chat.UnexpectedMessageType
}

func (ch *ChatHandler) deleteMessage(ctx context.Context, msg models.MessageInput, user int64, event string) error {
	message := models.Message{
		MessageID: msg.MessageID,
		EventID:   event,
		Type:      chat.DeleteMessage,
	}

	err := ch.useCase.DeleteMessage(ctx, msg.MessageID)
	if err != nil && !errors.Is(err, chat.ErrNotAllowedToDelete) {
		ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to delete message")
		return errors.Wrap(err, "failed to delete message")
	}
	if errors.Is(err, chat.ErrNotAllowedToDelete) {
		message.Error = chat.ErrNotAllowedToDelete.Error()
	}

	jsonMessage, err := json.Marshal(&message)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal json message")
	}

	err = ch.writeMessageToClients(ctx, jsonMessage, user, event)
	if err != nil {
		return errors.Wrap(err, "failed to write message to clients")
	}

	return nil
}

func (ch *ChatHandler) sendTextMessage(ctx context.Context, msg models.MessageInput, user int64, event string) error {
	messageID, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "failed to generate message id")
	}

	message := models.Message{
		MessageID:   messageID.String(),
		UserID:      user,
		EventID:     event,
		Text:        msg.Value,
		TimeCreated: time.Now().Unix(),
	}

	err = ch.useCase.CreateMessage(ctx, message)
	if err != nil {
		ch.logger.WithCtx(ctx).WithError(err).Error("failed to create message")
		return errors.Wrap(err, "failed to create message")
	}

	jsonMessage, err := json.Marshal(&message)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal json message")
	}

	err = ch.writeMessageToClients(ctx, jsonMessage, user, event)
	if err != nil {
		return errors.Wrap(err, "failed to write message to clients")
	}

	return nil
}

func (ch *ChatHandler) writeMessageToClients(ctx context.Context, jsonMessage []byte, user int64, event string) error {
	for _, client := range ch.messenger.Chats[event].Clients {
		err := client.Socket.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			ch.logger.WithCtx(ctx).WithError(err).Errorf("failed to write message to %d", client.UserID)
			ch.messenger.RemoveClientFromChat(event, user)
			continue
		}

		err = ch.useCase.UpdateLastCheckTime(ctx, event, client.UserID, time.Now().Unix())
		if err != nil {
			return errors.Wrap(err, "failed to update last check time")
		}
	}

	return nil
}
