package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"strconv"
)

func (gu *groupUseCase) GetAdministeredGroupByUserId(ctx context.Context, userId string) ([]models.Group, error) {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		gu.logger.WithError(err).Error("[GetAdministeredGroupByUserId] bad user id")
		return make([]models.Group, 0), err
	}

	groups, err := gu.repository.GetAdministeredGroupByUserId(ctx, userIdInt)
	if groups == nil {
		groups = make([]models.Group, 0)
	}

	return groups, err
}

func (gu *groupUseCase) CheckIfAdmin(ctx context.Context, userId string, groupId int) (bool, error) {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		gu.logger.WithError(err).Error("[GetAdministeredGroupByUserId] bad user id")
		return false, err
	}

	return gu.repository.CheckIfAdmin(ctx, userIdInt, groupId)
}
