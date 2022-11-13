package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/blacklist"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

var ErrBlacklistedEvent = errors.New("event data is blacklisted")

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

func (uc eventUsecase) UnsubscribeEvent(ctx context.Context, event string, user int64) error {
	err := uc.Events.UnSubscribe(ctx, event, user)
	if err != nil {
		return errors.Wrapf(err, "failed to unsubscribe event %s", event)
	}

	return nil
}

func (uc eventUsecase) Delete(ctx context.Context, event string, groupInfo models.GroupInfo) error {
	err := uc.Events.Delete(ctx, event, groupInfo)
	if err != nil {
		return errors.Wrapf(err, "failed to delete event %s", event)
	}

	return nil
}

func (uc eventUsecase) Change(ctx context.Context, event models.Event) error {
	err := uc.validateEvent(event)
	if err != nil {
		return errors.Wrap(err, " failed to validate event")
	}

	err = uc.Events.Update(ctx, event)
	if err != nil {
		return errors.Wrapf(err, "failed to update event %s", event.Uid)
	}

	return nil
}

func (uc eventUsecase) validateEvent(event models.Event) error {
	err := uc.BlackLister.Validate(blacklist.RowData{CheckData: event.Title})
	if err != nil {
		return errors.Wrap(ErrBlacklistedEvent, " failed to validate event's title")
	}

	err = uc.BlackLister.Validate(blacklist.RowData{CheckData: event.Description})
	if err != nil {
		return errors.Wrap(ErrBlacklistedEvent, " failed to validate event's description")
	}

	return nil
}
