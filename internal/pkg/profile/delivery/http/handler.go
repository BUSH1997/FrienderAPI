package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ProfileHandler struct {
	useCase profile.UseCase
	logger  *logrus.Logger
}

func NewProfileHandler(useCase profile.UseCase, logger *logrus.Logger) *ProfileHandler {
	return &ProfileHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (eh *ProfileHandler) GetOneProfile(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString == "" {
		eh.logger.Errorf("[GetOneProfile] failed get id")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		eh.logger.WithError(errors.Wrap(err, "failed to parse user id")).
			Errorf("failed to get user events")

		return ctx.NoContent(http.StatusInternalServerError)
	}

	profile, err := eh.useCase.GetOneProfile(ctx.Request().Context(), int(id))
	if err != nil {
		eh.logger.WithError(err).Errorf("[GetOneProfile] failed get one profile")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, profile)
}

func (eh *ProfileHandler) GetAllStatusesUser(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString == "" {
		eh.logger.Errorf("[GetAllStatusesUser] failed get id")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		eh.logger.WithError(errors.Wrap(err, "failed to parse user id")).
			Errorf("failed to get user events")

		return ctx.NoContent(http.StatusInternalServerError)
	}

	statuses, err := eh.useCase.GetAllProfileStatuses(ctx.Request().Context(), int(id))
	if err != nil {
		eh.logger.WithError(err).Errorf("[GetAllStatusesUser] failed get all statuses user")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, statuses)
}

func (eh *ProfileHandler) ChangeProfile(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString == "" {
		eh.logger.Errorf("[ChangeProfile] failed get id")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		eh.logger.WithError(errors.Wrap(err, "failed to parse user id")).
			Errorf("failed to get user events")

		return ctx.NoContent(http.StatusInternalServerError)
	}

	var newProfileData models.ChangeProfile
	if err := ctx.Bind(&newProfileData); err != nil {
		eh.logger.WithError(err).Errorf("[ChangeProfile] failed bind change profile")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	newProfileData.ProfileId = id

	if err := eh.useCase.UpdateProfile(ctx.Request().Context(), newProfileData); err != nil {
		eh.logger.WithError(err).Errorf("[ChangeProfile] failed change profile")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)
}

func (eh *ProfileHandler) ChangePriorityEvent(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString != "" {
		eh.logger.Errorf("[ChangePriorityEvent] failed get id")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		eh.logger.WithError(errors.Wrap(err, "failed to parse user id")).
			Errorf("failed to get user events")

		return ctx.NoContent(http.StatusInternalServerError)
	}

	var newPriorityEvent models.UidEventPriority
	if err := ctx.Bind(&newPriorityEvent); err != nil {
		eh.logger.WithError(err).Errorf("[ChangePriorityEvent] failed bind priority event")
		return ctx.NoContent(http.StatusInternalServerError)
	}
	newPriorityEvent.UidUser = int(id)

	if err := eh.useCase.ChangeEventPriority(ctx.Request().Context(), newPriorityEvent); err != nil {
		eh.logger.WithError(err).Errorf("[ChangePriorityEvent] failed change priority event")
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)
}
