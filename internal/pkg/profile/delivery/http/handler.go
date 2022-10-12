package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
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
		return ctx.NoContent(http.StatusInternalServerError)
	}

	profile, err := eh.useCase.GetOneProfile(idString)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, profile)
}

func (eh *ProfileHandler) GetAllStatusesUser(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString == "" {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	statuses, err := eh.useCase.GetAllStatusesUser(idString)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, statuses)
}

func (eh *ProfileHandler) ChangeProfile(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString == "" {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	var newProfileData models.ChangeProfile
	if err := ctx.Bind(&newProfileData); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	newProfileData.ProfileId = idString

	if err := eh.useCase.ChangeProfile(newProfileData); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)
}

func (eh *ProfileHandler) ChangePriorityEvent(ctx echo.Context) error {
	idString := ctx.Param("id")
	if idString != "" {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	var newPriorityEvent models.UidEventPriority
	if err := ctx.Bind(&newPriorityEvent); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	newPriorityEvent.UidUser = idString

	if err := eh.useCase.ChangePriorityEvent(newPriorityEvent); err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)
}
