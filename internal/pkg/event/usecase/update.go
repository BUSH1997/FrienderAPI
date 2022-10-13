package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (u eventUsecase) Update(ctx context.Context, event models.Event) error {
	err := u.Events.Update(ctx, event)
	if err != nil {
		return errors.Wrap(err, "failed to update public event in usecase")
	}

	return nil
}

func (u eventUsecase) SubscribeEvent(ctx context.Context, id models.UserIdEventId) error {
	return nil
}

func (u eventUsecase) UnsubscribeEvent(ctx context.Context, id models.UserIdEventId) error {
	return nil
}

func (u eventUsecase) DeleteEvent(ctx context.Context, id models.UserIdEventId) error {
	return nil
}

func (u eventUsecase) ChangeEvent(ctx context.Context, event models.Event) error {
	return nil
}
