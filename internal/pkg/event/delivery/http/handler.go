package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/queryParamParser"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
	queryParam := ctx.QueryParams()
	filter, err := queryParamParser.ParseGetAllEvents(queryParam)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	events, err := eh.useCase.GetAll(ctx.Request().Context(), filter)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to get events")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, events)
}

func (eh *EventHandler) GetEventsUser(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		eh.logger.WithError(errors.Wrap(err, "failed to parse user id")).
			Errorf("failed to get user events")

		return ctx.NoContent(http.StatusInternalServerError)
	}

	events, err := eh.useCase.GetUserEvents(ctx.Request().Context(), id)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to get user events")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, events)
}

func (eh *EventHandler) Get(ctx echo.Context) error {
	eventParams := models.GetEventParams{}

	idString := ctx.QueryParam("id")
	if idString == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}
	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		eh.logger.WithError(errors.Wrap(err, "failed to parse user id")).
			Errorf("failed to get user events")

		return ctx.NoContent(http.StatusBadRequest)
	}

	eventParams.UserID = id

	isSubString := ctx.QueryParam("is_sub")
	if isSubString != "" {
		isSub, err := strconv.ParseBool(isSubString)
		if err != nil {
			eh.logger.WithError(errors.Wrap(err, "failed to parse sub param")).
				Errorf("failed to get user events")

			return ctx.NoContent(http.StatusBadRequest)
		}

		eventParams.IsSubscriber = models.DefinedBool(isSub)
	}

	isActiveString := ctx.QueryParam("is_active")
	if isActiveString != "" {
		isActive, err := strconv.ParseBool(isActiveString)
		if err != nil {
			eh.logger.WithError(errors.Wrap(err, "failed to parse active param")).
				Errorf("failed to get user events")

			return ctx.NoContent(http.StatusBadRequest)
		}

		eventParams.IsActive = models.DefinedBool(isActive)
	}

	isOwnerString := ctx.QueryParam("is_active")
	if isOwnerString != "" {
		isOwner, err := strconv.ParseBool(isOwnerString)
		if err != nil {
			eh.logger.WithError(errors.Wrap(err, "failed to parse owner param")).
				Errorf("failed to get user events")

			return ctx.NoContent(http.StatusBadRequest)
		}

		eventParams.IsOwner = models.DefinedBool(isOwner)
	}

	events, err := eh.useCase.Get(ctx.Request().Context(), eventParams)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to get user events")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, events)
}

func (eh *EventHandler) SubscribeEvent(ctx echo.Context) error {
	eventID := ctx.Param("id")
	if eventID == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var userId models.UserId
	if err := ctx.Bind(&userId); err != nil {
		eh.logger.WithError(err).Errorf("failed to subscribe event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err := eh.useCase.SubscribeEvent(ctx.Request().Context(), int64(userId.Id), eventID)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to subscribe event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.NoContent(http.StatusInternalServerError)
}

func (eh *EventHandler) UnsubscribeEvent(ctx echo.Context) error {
	eventID := ctx.Param("id")
	if eventID == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var userId models.UserId
	if err := ctx.Bind(&userId); err != nil {
		eh.logger.WithError(err).Errorf("failed to unsubscribe event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err := eh.useCase.UnsubscribeEvent(ctx.Request().Context(), int64(userId.Id), eventID)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to unsubscribe event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.NoContent(http.StatusInternalServerError)
}

func (eh *EventHandler) DeleteEvent(ctx echo.Context) error {
	eventID := ctx.Param("id")
	if eventID == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	var userId models.UserId
	if err := ctx.Bind(&userId); err != nil {
		eh.logger.WithError(err).Errorf("failed to delete event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err := eh.useCase.DeleteEvent(ctx.Request().Context(), int64(userId.Id), eventID)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to delete event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.NoContent(http.StatusInternalServerError)
}

func (eh *EventHandler) ChangeEvent(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	_, err := strconv.Atoi(idString)
	if err != nil {
		eh.logger.WithError(errors.Wrap(err, "failed to parse user id")).
			Errorf("failed to delete event")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	var event models.Event
	if err := ctx.Bind(&event); err != nil {
		eh.logger.WithError(err).Errorf("failed to delete event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err = eh.useCase.ChangeEvent(ctx.Request().Context(), event)
	if err != nil {
		eh.logger.WithError(err).Errorf("failed to delete event")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	return ctx.NoContent(http.StatusInternalServerError)
}

func (eh *EventHandler) GetAllCategory(ctx echo.Context) error {
	categories, err := eh.useCase.GetAllCategories(ctx.Request().Context())
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, categories)
}
