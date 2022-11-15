package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

func (gu *groupUseCase) Create(ctx context.Context, group models.GroupInput) error {
	ctx = gu.logger.WithCaller(ctx)

	if err := gu.repository.Create(ctx, group); err != nil {
		gu.logger.WithError(err).Errorf("[Create] use case")
		return err
	}

	if err := gu.repositoryUser.Create(ctx, int64(group.GroupId), true); err != nil {
		gu.logger.WithError(err).Errorf("[Create] use case")
		return err
	}

	return nil
}
