package http

import (
	"fmt"
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
	idString := ctx.Param("id")
	if idString == "" {
		fmt.Println("LOL3")
		return ctx.NoContent(http.StatusBadRequest)
	}

	event, err := eh.useCase.GetEventById(ctx.Request().Context(), idString)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, event)
}

func (eh *EventHandler) GetEvents(ctx echo.Context) error {
	events, err := eh.useCase.GetAll(ctx.Request().Context())
	if err != nil {
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
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, events)
}
