package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
	"strconv"
)

func (gu *groupUseCase) GetAdministeredGroupByUserId(ctx context.Context, userId string) ([]models.Group, error) {
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
	group, err := gu.repository.Get(ctx, userID)
	if err != nil {
		gu.logger.WithError(err).Error("failed to get group")
		return models.Group{}, errors.Wrap(err, "failed to get group in usecase")
	}

	return group, err
}

func (gu *groupUseCase) CheckIfAdmin(ctx context.Context, userId string, groupId int64) (bool, error) {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		gu.logger.WithError(err).Error("[GetAdministeredGroupByUserId] bad user id")
		return false, err
	}

	return gu.repository.CheckIfAdmin(ctx, userIdInt, groupId)
}
