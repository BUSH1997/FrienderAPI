package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/usecase"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type EventHandler struct {
	useCase event.Usecase
	logger  *logrus.Logger
}

func NewEventHandler(usecase event.Usecase, logger *logrus.Logger) *EventHandler {
	return &EventHandler{
		useCase: usecase,
		logger:  logger,
	}
}

func (eh *EventHandler) Create(ctx echo.Context) error {
	var newEvent models.Event
	if err := ctx.Bind(&newEvent); err != nil {
		eh.logger.WithError(err).Errorf("failed to create event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	event, err := eh.useCase.Create(ctx.Request().Context(), newEvent)
	if errors.Is(err, usecase.ErrBlacklistedEvent) {
		eh.logger.WithError(err).Errorf("failed to create event")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to create event")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, event)
}

func (eh *EventHandler) GetOneEvent(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString == "" {
		err := errors.New("event id is empty")
		eh.logger.WithError(err).Errorf("failed to get one event")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	event, err := eh.useCase.GetEventById(ctx.Request().Context(), idString)
	if err != nil {
		eh.logger.WithField("event", idString).WithError(err).Errorf("failed to get one event")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, event)
}

func (eh *EventHandler) Get(ctx echo.Context) error {
	eventParams := models.GetEventParams{}
	if err := ctx.Bind(&eventParams); err != nil {
		eh.logger.WithError(err).Errorf("failed to bind events get params")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// eventParams.UserID = contextlib.GetUser(ctx.Request().Context())
	events, err := eh.useCase.Get(ctx.Request().Context(), eventParams)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to get user events")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, events)
}

func (eh *EventHandler) SubscribeEvent(ctx echo.Context) error {
	eventID := ctx.Param("id")
	if eventID == "" {
		err := errors.New("event id is empty")
		eh.logger.WithError(err).Errorf("failed to subscribe event")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := eh.useCase.SubscribeEvent(ctx.Request().Context(), eventID)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to subscribe event")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, "successfully subscribed event")
}

func (eh *EventHandler) UnsubscribeEvent(ctx echo.Context) error {
	eventID := ctx.Param("id")
	if eventID == "" {
		err := errors.New("event id is empty")
		eh.logger.WithError(err).Errorf("failed to unsubscribe event")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := eh.useCase.UnsubscribeEvent(ctx.Request().Context(), eventID)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to unsubscribe event")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "successfully unsubscribed event")
}

func (eh *EventHandler) DeleteEvent(ctx echo.Context) error {
	eventID := ctx.Param("id")
	if eventID == "" {
		err := errors.New("event id is empty")
		eh.logger.WithError(err).Errorf("failed to delete event")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	var groupInfo models.GroupInfo
	if err := ctx.Bind(&groupInfo); err != nil {
		eh.logger.WithError(err).Errorf("failed to bind group info event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err := eh.useCase.Delete(ctx.Request().Context(), eventID, groupInfo)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to delete event")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "successfully deleted event")
}

func (eh *EventHandler) ChangeEvent(ctx echo.Context) error {
	var event models.Event
	if err := ctx.Bind(&event); err != nil {
		eh.logger.WithError(err).Errorf("failed to bind event")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := eh.useCase.Change(ctx.Request().Context(), event)
	if errors.Is(err, usecase.ErrBlacklistedEvent) {
		eh.logger.WithError(err).Errorf("failed to create event")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to change event")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, event)
}

func (eh *EventHandler) GetAllCategory(ctx echo.Context) error {
	categories, err := eh.useCase.GetAllCategories(ctx.Request().Context())
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to get all categories")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, categories)
}

func (eh *EventHandler) UpdateAlbum(ctx echo.Context) error {
	var updateAlbumInfo models.UpdateAlbumInfo
	if err := ctx.Bind(&updateAlbumInfo); err != nil {
		eh.logger.WithError(err).Errorf("failed to bind UpdateAlbumInfo")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := eh.useCase.UpdateAlbum(ctx.Request().Context(), updateAlbumInfo); err != nil {
		eh.logger.WithError(err).Errorf("failed to updateAlbum")
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, updateAlbumInfo)
}
