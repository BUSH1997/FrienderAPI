package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"strconv"
)

func (gu *groupUseCase) GetAdministeredGroupByUserId(ctx context.Context, userId string) ([]models.Group, error) {
	ctx = gu.logger.WithCaller(ctx)

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		gu.logger.WithError(err).Error("[GetAdministeredGroupByUserId] bad user id")
		return make([]models.Group, 0), errors.Wrap(err, "failed to parse user id from string")
	}

	groups, err := gu.repository.GetAdministeredGroupByUserId(ctx, userIdInt)
	if err != nil {
		gu.logger.WithError(err).Error("[GetAdministeredGroupByUserId] bad user id")
		return make([]models.Group, 0), errors.Wrap(err, "failed to get administrated groups")
	}
	if groups == nil {
		groups = make([]models.Group, 0)
	}

	return groups, err
}

func (gu *groupUseCase) Get(ctx context.Context, userID int64) (models.Group, error) {
	ctx = gu.logger.WithCaller(ctx)

	group, err := gu.repository.Get(ctx, userID)
	if err != nil {
		gu.logger.WithError(err).Error("failed to get group")
		return models.Group{}, errors.Wrap(err, "failed to get group in usecase")
	}

	return group, err
}

func (gu *groupUseCase) CheckIfAdmin(ctx context.Context, userId int64, groupId int64) (bool, error) {
	ctx = gu.logger.WithCaller(ctx)

	isAdmin, err := gu.repository.CheckIfAdmin(ctx, userId, groupId)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if admin in usecase")
	}

	return isAdmin, nil
}
