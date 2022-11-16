package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/api/errors/convert"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ProfileHandler struct {
	useCase profile.UseCase
	logger  hardlogger.Logger
}

func NewProfileHandler(useCase profile.UseCase, logger hardlogger.Logger) *ProfileHandler {
	return &ProfileHandler{
		useCase: useCase,
		logger:  logger,
	}
}

func (eh *ProfileHandler) GetOneProfile(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	idString := echoCtx.Param("id")
	if idString == "" {
		err := errors.New("empty profile id")
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed get user id param")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	id, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to parse user id %s", id)
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	profile, err := eh.useCase.GetOneProfile(ctx, id)
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed get profile %d", id)
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, profile)
}

func (eh *ProfileHandler) GetAllStatusesUser(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	statuses, err := eh.useCase.GetAllProfileStatuses(ctx)
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to get all user statuses")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, statuses)
}

func (eh *ProfileHandler) ChangeProfile(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	var newProfileData models.ChangeProfile
	if err := echoCtx.Bind(&newProfileData); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to bind change profile data")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	newProfileData.ProfileId = contextlib.GetUser(ctx)

	if err := eh.useCase.UpdateProfile(ctx, newProfileData); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to change profile")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.NoContent(http.StatusOK)
}

func (eh *ProfileHandler) ChangePriorityEvent(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	var newPriorityEvent models.UidEventPriority
	if err := echoCtx.Bind(&newPriorityEvent); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed bind priority event")
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}
	newPriorityEvent.UidUser = int(contextlib.GetUser(ctx))

	if err := eh.useCase.ChangeEventPriority(ctx, newPriorityEvent); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed change event priority")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.NoContent(http.StatusOK)
}

func (eh *ProfileHandler) Subscribe(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	var profileForSubscribeID models.UserId
	if err := echoCtx.Bind(&profileForSubscribeID); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed bind subscribe data")
		return echoCtx.JSON(http.StatusBadRequest, err)
	}

	userID := contextlib.GetUser(ctx)
	err := eh.useCase.Subscribe(ctx, userID, int64(profileForSubscribeID.Id))
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to subscribe to %d", profileForSubscribeID.Id)
		return echoCtx.JSON(http.StatusBadRequest, err)
	}

	return echoCtx.JSON(http.StatusOK, profileForSubscribeID)
}

func (eh *ProfileHandler) UnSubscribe(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	var profileForUnsubscribeId models.UserId
	if err := echoCtx.Bind(&profileForUnsubscribeId); err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Error("failed bind unsubscribe data")
		return echoCtx.JSON(http.StatusBadRequest, err)
	}

	userID := contextlib.GetUser(ctx)

	err := eh.useCase.UnSubscribe(ctx, userID, int64(profileForUnsubscribeId.Id))
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).
			Errorf("failed to unsubscribe from %d", profileForUnsubscribeId.Id)
		return echoCtx.JSON(http.StatusBadRequest, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, profileForUnsubscribeId)
}

func (eh *ProfileHandler) GetSubscribe(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	userID := contextlib.GetUser(ctx)

	subscribe, err := eh.useCase.GetSubscribe(ctx, userID)
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to get subscriptions")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, subscribe)
}

func (eh *ProfileHandler) GetFriends(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	userID := contextlib.GetUser(ctx)

	friends, err := eh.useCase.GetFriends(ctx, userID)
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Error("failed to get friends")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, friends)
}

func (eh *ProfileHandler) GetAllCities(echoCtx echo.Context) error {
	ctx := eh.logger.WithCaller(echoCtx.Request().Context())

	cities, err := eh.useCase.GetCities(ctx)
	if err != nil {
		eh.logger.WithCtx(ctx).WithError(err).Errorf("failed to get all cities")
		return echoCtx.JSON(http.StatusInternalServerError, convert.DeliveryError(err).Error())
	}

	return echoCtx.JSON(http.StatusOK, cities)
}
