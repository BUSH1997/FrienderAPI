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

func (gh *GroupHandler) Update(ctx echo.Context) error {
	var newGroupData models.Group
	if err := ctx.Bind(&newGroupData); err != nil {
		gh.logger.WithError(err).Errorf("failed to bind group data")
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err := gh.useCase.Update(ctx.Request().Context(), newGroupData)
	if err != nil {
		gh.logger.WithError(err).Errorf("failed to update group data")
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, newGroupData)
}

func (gh *GroupHandler) GetAdministeredGroup(ctx echo.Context) error {
	userId := ctx.QueryParam("user_id")
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

func (gh *GroupHandler) Get(ctx echo.Context) error {
	groupIdString := ctx.QueryParam("group_id")
	if groupIdString == "" {
		gh.logger.Error("[Get], bad group_id")
		return ctx.JSON(http.StatusBadRequest, errors.New(" bad group_id").Error())
	}

	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		gh.logger.WithError(err).Error("failed to parse user id from string")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	group, err := gh.useCase.Get(ctx.Request().Context(), int64(groupId))
	if err != nil {
		gh.logger.WithError(err).Errorf("failed to get group by id %d", groupId)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, group)
}

func (gh *GroupHandler) IsAdmin(ctx echo.Context) error {
	userId := ctx.Request().Header.Get("X-User-ID")
	if userId == "" {
		gh.logger.Error("[IsAdmin], bad x-user-id")
		return ctx.JSON(http.StatusBadRequest, errors.New("bad x-user-id"))
	}

	groupIdString := ctx.QueryParam("group_id")
	if groupIdString == "" {
		gh.logger.Error("[IsAdmin], bad group_id")
		return ctx.JSON(http.StatusBadRequest, errors.New("bad group_id"))
	}

	groupId, err := strconv.ParseInt(groupIdString, 10, 32)
	if err != nil {
		gh.logger.WithError(errors.Wrap(err, "failed to parse owner param")).
			Errorf("failed to get group id")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	isAdmin, err := gh.useCase.CheckIfAdmin(ctx.Request().Context(), userId, groupId)
	if err != nil {
		gh.logger.WithError(err).Errorf("[IsAdmin], error in useCase")
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, isAdmin)
}
