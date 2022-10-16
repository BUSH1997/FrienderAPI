package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (uc eventUsecase) Update(ctx context.Context, event models.Event) error {
	err := uc.Events.Update(ctx, event)
	if err != nil {
		return errors.Wrap(err, "failed to update public event in usecase")
	}
	return nil
}

func (uc eventUsecase) SubscribeEvent(ctx context.Context, user int64, event string) error {
	err := uc.Events.Subscribe(ctx, user, event)
	if err != nil {
		return errors.Wrapf(err, "failed to subscribe event %s", event)
	}

	return nil
}

func (uc eventUsecase) UnsubscribeEvent(ctx context.Context, user int64, event string) error {
	err := uc.Events.UnSubscribe(ctx, user, event)
	if err != nil {
		return errors.Wrapf(err, "failed to unsubscribe event %s", event)
	}

	return nil
}

func (uc eventUsecase) DeleteEvent(ctx context.Context, user int64, event string) error {
	return nil
}

func (uc eventUsecase) ChangeEvent(ctx context.Context, event models.Event) error {
	return nil
}
