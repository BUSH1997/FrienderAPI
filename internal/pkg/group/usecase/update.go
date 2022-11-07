package usecase

import (
	"context"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
	"strconv"
)

func (gu *groupUseCase) Update(ctx context.Context, group models.GroupInput) error {
	err := gu.repository.Update(ctx, group)
	if err != nil {
		return errors.Wrap(err, "failed to update group in usecase")
	}

	return nil
}

func (gu *groupUseCase) ApproveEvent(ctx context.Context, approveInfo models.ApproveEvent) error {
	userId := strconv.Itoa(int(contextlib.GetUser(ctx)))
	isAdmin, err := gu.CheckIfAdmin(ctx, userId, int64(approveInfo.GroupId))
	if err != nil {
		gu.logger.WithError(err).Errorf("[ApproveEvent] faile checkIfAdmin")
		return err
	}
	if !isAdmin {
		return errors.New("Try approve event not admin")
	}

	err = gu.repository.ApproveEvent(ctx, approveInfo)
	if err != nil {
		gu.logger.WithError(err).Errorf("[ApproveEvent] faile approve event")
		return err
	}

	return nil
}
