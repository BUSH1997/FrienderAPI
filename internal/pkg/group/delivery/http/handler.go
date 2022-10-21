package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/group"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type GroupHandler struct {
	logger  *logrus.Logger
	useCase group.UseCase
}

func New(logger *logrus.Logger, useCase group.UseCase) *GroupHandler {
	return &GroupHandler{
		logger:  logger,
		useCase: useCase,
	}
}

func (gh *GroupHandler) CreateGroup(ctx echo.Context) error {
	var newGroup models.Group
	userId := ctx.Request().Header.Get("X-User-ID")
	if userId == "" {
		gh.logger.Error("[GetAdministeredGroup], bad x-user-id")
		return ctx.JSON(http.StatusBadRequest, errors.New("bad x-user-id"))
	}

	if err := ctx.Bind(&newGroup); err != nil {
		gh.logger.WithError(err).Errorf("[Create group], error json")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		gh.logger.WithError(err).Error("[CreateGroup] bad user id")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if userIdInt != newGroup.UserId {
		gh.logger.Error("[CreateGroup] try to create group not with your id")
		return ctx.JSON(http.StatusBadRequest, errors.New("try to create group not with your id"))
	}

	err = gh.useCase.Create(ctx.Request().Context(), newGroup)
	if err != nil {
		gh.logger.WithError(err).Errorf("[Create group], error in useCase")
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, newGroup)
}

func (gh *GroupHandler) GetAdministeredGroup(ctx echo.Context) error {
	userId := ctx.Request().Header.Get("X-User-ID")
	if userId == "" {
		gh.logger.Error("[GetAdministeredGroup], bad x-user-id")
		return ctx.JSON(http.StatusBadRequest, errors.New("bad x-user-id"))
	}

	groups, err := gh.useCase.GetAdministeredGroupByUserId(ctx.Request().Context(), userId)
	if err != nil {
		gh.logger.WithError(err).Errorf("[GetAdministeredGroup], error in useCase")
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, groups)
}
