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

func (uc eventUsecase) SubscribeEvent(ctx context.Context, event string) error {
	err := uc.Events.Subscribe(ctx, event)
	if err != nil {
		return errors.Wrapf(err, "failed to subscribe event %s", event)
	}

	return nil
}

func (uc eventUsecase) UnsubscribeEvent(ctx context.Context, event string) error {
	err := uc.Events.UnSubscribe(ctx, event)
	if err != nil {
		return errors.Wrapf(err, "failed to unsubscribe event %s", event)
	}

	return nil
}

func (uc eventUsecase) Delete(ctx context.Context, event string) error {
	err := uc.Events.Delete(ctx, event)
	if err != nil {
		return errors.Wrapf(err, "failed to delete event %s", event)
	}

	return nil
}

func (uc eventUsecase) Change(ctx context.Context, event models.Event) error {
	err := uc.Events.Update(ctx, event)
	if err != nil {
		return errors.Wrapf(err, "failed to update event %s", event.Uid)
	}

	return nil
}
