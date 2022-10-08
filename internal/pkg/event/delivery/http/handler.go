package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
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

func (eh *EventHandler) CreateEvent(ctx echo.Context) error {
	var newEvent models.Event
	if err := ctx.Bind(&newEvent); err != nil {
		eh.logger.WithError(err).Errorf("failed to create event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	event, err := eh.useCase.Create(ctx.Request().Context(), newEvent)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to create event")
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, event)
}

func (eh *EventHandler) GetOneEvent(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString == "" {
		eh.logger.WithError(errors.New("event id is empty")).Errorf("failed to get one event")
		return ctx.NoContent(http.StatusBadRequest)
	}

	event, err := eh.useCase.GetEventById(ctx.Request().Context(), idString)
	if err != nil {
		eh.logger.WithField("event", idString).WithError(err).Errorf("failed to get one event")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, event)
}

func (eh *EventHandler) GetEvents(ctx echo.Context) error {
	events, err := eh.useCase.GetAll(ctx.Request().Context())
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to get events")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, events)
}

func (eh *EventHandler) GetEventsUser(ctx echo.Context) error {
	idString := ctx.QueryParams().Get("id")
	if idString == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	events, err := eh.useCase.GetUserEvents(ctx.Request().Context(), idString)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to get user events")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, events)
}
