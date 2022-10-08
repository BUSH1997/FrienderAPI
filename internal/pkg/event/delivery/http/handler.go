package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EventHandler struct {
	useCase event.Usecase
}

func NewEventHandler(usecase event.Usecase) *EventHandler {
	return &EventHandler{
		useCase: usecase,
	}
}

func (eh *EventHandler) CreateEvent(ctx echo.Context) error {
	var newEvent models.Event
	if err := ctx.Bind(&newEvent); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if err := eh.useCase.Create(ctx.Request().Context(), newEvent); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, 5)
}

func (eh *EventHandler) GetOneEvent(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, 5)
}

func (eh *EventHandler) GetEvents(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, 5)
}

func (eh *EventHandler) GetEventsUser(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, 5)
}
