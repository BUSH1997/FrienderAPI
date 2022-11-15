package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/group"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type GroupHandler struct {
	logger  hardlogger.Logger
	useCase group.UseCase
}

func New(logger hardlogger.Logger, useCase group.UseCase) *GroupHandler {
	return &GroupHandler{
		logger:  logger,
		useCase: useCase,
	}
}

func (gh *GroupHandler) CreateGroup(echoCtx echo.Context) error {
	ctx := gh.logger.WithCaller(echoCtx.Request().Context())

	var newGroup models.GroupInput
	if err := echoCtx.Bind(&newGroup); err != nil {
		gh.logger.WithCtx(ctx).WithError(err).Errorf("[Create group], error json")
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	userID := context.GetUser(ctx)
	if int(userID) != newGroup.UserId {
		err := errors.New("attempt to create group with alien id")
		gh.logger.WithCtx(ctx).WithError(err).Error("failed to try to create group not with your id")
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	err := gh.useCase.Create(ctx, newGroup)
	if err != nil {
		gh.logger.WithCtx(ctx).WithError(err).Errorf("failed to create group")
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, newGroup)
}

func (gh *GroupHandler) Update(echoCtx echo.Context) error {
	ctx := gh.logger.WithCaller(echoCtx.Request().Context())

	var newGroupData models.GroupInput
	if err := echoCtx.Bind(&newGroupData); err != nil {
		gh.logger.WithCtx(ctx).WithError(err).Errorf("failed to bind group data")
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	err := gh.useCase.Update(ctx, newGroupData)
	if err != nil {
		gh.logger.WithCtx(ctx).WithError(err).Errorf("failed to update group data")
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, newGroupData)
}

func (gh *GroupHandler) GetAdministeredGroup(echoCtx echo.Context) error {
	ctx := gh.logger.WithCaller(echoCtx.Request().Context())

	userId := echoCtx.QueryParam("user_id")
	if userId == "" {
		err := errors.New("empty user id")
		gh.logger.WithCtx(ctx).WithError(err).Error("failed to parse user id param")
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	groups, err := gh.useCase.GetAdministeredGroupByUserId(ctx, userId)
	if err != nil {
		gh.logger.WithCtx(ctx).WithError(err).
			Errorf("failed to get administrated group by user id %d", userId)
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, groups)
}

func (gh *GroupHandler) Get(echoCtx echo.Context) error {
	ctx := gh.logger.WithCaller(echoCtx.Request().Context())

	groupIdString := echoCtx.QueryParam("group_id")
	if groupIdString == "" {
		err := errors.New("empty user id")
		gh.logger.WithCtx(ctx).WithError(err).Error("failed to parse group id param")
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		gh.logger.WithCtx(ctx).WithError(err).Error("failed to parse group id %d from string", groupIdString)
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	group, err := gh.useCase.Get(ctx, int64(groupId))
	if err != nil {
		gh.logger.WithCtx(ctx).WithError(err).Errorf("failed to get group by id %d", groupId)
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, group)
}

func (gh *GroupHandler) IsAdmin(echoCtx echo.Context) error {
	ctx := gh.logger.WithCaller(echoCtx.Request().Context())

	userId := context.GetUser(ctx)

	groupIdString := echoCtx.QueryParam("group_id")
	if groupIdString == "" {
		err := errors.New("empty group id")
		gh.logger.WithCtx(ctx).WithError(err).Error("failed to parse group id param")
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	groupId, err := strconv.ParseInt(groupIdString, 10, 32)
	if err != nil {
		gh.logger.WithCtx(ctx).WithError(errors.Wrap(err, "failed to parse owner param")).
			Errorf("failed to get group id")
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	isAdmin, err := gh.useCase.CheckIfAdmin(ctx, userId, groupId)
	if err != nil {
		gh.logger.WithCtx(ctx).WithError(err).Errorf("failed to check if admin")
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, isAdmin)
}

func (gh *GroupHandler) ApproveEvent(echoCtx echo.Context) error {
	ctx := gh.logger.WithCaller(echoCtx.Request().Context())

	var approveEvent models.ApproveEvent
	if err := echoCtx.Bind(&approveEvent); err != nil {
		gh.logger.WithCtx(ctx).WithError(err).Errorf("failed to bind approve event data")
		return echoCtx.JSON(http.StatusBadRequest, err.Error())
	}

	err := gh.useCase.ApproveEvent(ctx, approveEvent)
	if err != nil {
		gh.logger.WithCtx(ctx).WithError(err).Errorf("failed to approve event %d", approveEvent.EventUid)
		return echoCtx.JSON(http.StatusInternalServerError, err.Error())
	}

	return echoCtx.JSON(http.StatusOK, approveEvent)
}
