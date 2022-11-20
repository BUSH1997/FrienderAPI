package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/api/errors/convert"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/usecase"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type EventHandler struct {
	useCase event.Usecase
	logger  hardlogger.Logger
}

func NewEventHandler(usecase event.Usecase, logger hardlogger.Logger) *EventHandler {
	return &EventHandler{
		useCase: usecase,
		logger:  logger,
	}
}

func (eh *EventHandler) Create(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	var newEvent models.Event
	if err := echoCtx.Bind(&newEvent); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to bind new event data")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	if err := echoCtx.Validate(&newEvent); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed validate event")
		return echoCtx.JSON(http.StatusBadRequest, errors.New("Failed validate data").Error())
	}

	event, err := eh.useCase.Create(ctx, newEvent)

	if errors.Is(err, usecase.ErrBlacklistedEvent) {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to create event")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to create event")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, event)
}

func (eh *EventHandler) GetOneEvent(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	idString := echoCtx.Param("id")
	if idString == "" {
		err := errors.New("event id is empty")
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to get one event")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	event, err := eh.useCase.GetEventById(ctx, idString)
	if err != nil {
		eh.logger.WithCtx(ctx).WithFields(hardlogger.Fields{
			"event": idString,
		}).WithError(err).Errorf("failed to get one event")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, event)
}

func (eh *EventHandler) Get(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	eventParams := models.GetEventParams{}
	if err := echoCtx.Bind(&eventParams); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to bind events get params")
		return echoCtx.JSON(http.StatusBadRequest, err)
	}

	events, err := eh.useCase.Get(ctx, eventParams)
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to get events")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, events)
}

func (eh *EventHandler) SubscribeEvent(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	eventID := echoCtx.Param("id")
	if eventID == "" {
		err := errors.New("event id is empty")
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to subscribe event")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	err := eh.useCase.SubscribeEvent(echoCtx.Request().Context(), eventID)
	if err != nil {
		eh.logger.WithCtx(ctx).WithFields(hardlogger.Fields{
			"event": eventID,
		}).WithError(err).Errorf("failed to subscribe event")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, "successfully subscribed event")
}

func (eh *EventHandler) UnsubscribeEvent(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	eventID := echoCtx.Param("id")
	if eventID == "" {
		err := errors.New("event id is empty")
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to unsubscribe event")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	input := models.UnsubscribeEventInput{}
	if err := echoCtx.Bind(&input); err != nil {
		eh.logger.WithError(err).Errorf("failed to bind unsubscribe event input")
		return echoCtx.JSON(http.StatusBadRequest, err)
	}

	user := context.GetUser(ctx)
	if input.User != 0 {
		user = input.User
	}

	err := eh.useCase.UnsubscribeEvent(ctx, eventID, user)
	if err != nil {
		eh.logger.WithCtx(ctx).WithFields(hardlogger.Fields{
			"event": eventID,
		}).WithError(err).Errorf("failed to unsubscribe event")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, "successfully unsubscribed event")
}

func (eh *EventHandler) DeleteEvent(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	eventID := echoCtx.Param("id")
	if eventID == "" {
		err := errors.New("event id is empty")
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to get event id param")

		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	var groupInfo models.GroupInfo
	if err := echoCtx.Bind(&groupInfo); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to bind group info event")
		return echoCtx.JSON(http.StatusBadRequest, err)
	}

	err := eh.useCase.Delete(ctx, eventID, groupInfo)
	if err != nil {
		eh.logger.WithCtx(ctx).WithFields(hardlogger.Fields{
			"event": eventID,
		}).WithError(err).Errorf("failed to delete event")

		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, "successfully deleted event")
}

func (eh *EventHandler) ChangeEvent(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	var event models.Event
	if err := echoCtx.Bind(&event); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to bind event")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	err := eh.useCase.Change(ctx, event)
	if errors.Is(err, usecase.ErrBlacklistedEvent) {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to create event")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to change event")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, event)
}

func (eh *EventHandler) GetAllCategory(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	categories, err := eh.useCase.GetAllCategories(ctx)
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to get all categories")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, categories)
}

func (eh *EventHandler) UpdateAlbum(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	var updateAlbumInfo models.UpdateAlbumInfo
	if err := echoCtx.Bind(&updateAlbumInfo); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to bind album info")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	userId := context.GetUser(ctx)
	updateAlbumInfo.UidAlbum = strconv.FormatInt(userId, 10) + "_" + updateAlbumInfo.UidAlbum

	if err := eh.useCase.UpdateAlbum(ctx, updateAlbumInfo); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to update album")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, updateAlbumInfo)
}
