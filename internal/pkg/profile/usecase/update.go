package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (uc *UseCase) UpdateProfile(ctx context.Context, profile models.ChangeProfile) error {
	ctx = uc.Logger.WithCaller(ctx)

	err := uc.profileRepository.UpdateProfile(ctx, profile)
	if err != nil {
		return errors.Wrap(err, "failed to update profile")
	}

	return nil
}

func (uc *UseCase) ChangeEventPriority(ctx context.Context, eventPriority models.UidEventPriority) error {
	ctx = uc.Logger.WithCaller(ctx)

	err := uc.eventRepository.UpdateEventPriority(ctx, eventPriority)
	if err != nil {
		return errors.Wrap(err, "failed to update event priority")
	}

	return nil
}

func (uc *UseCase) Subscribe(ctx context.Context, userId int64, groupId int64) error {
	ctx = uc.Logger.WithCaller(ctx)

	err := uc.profileRepository.Subscribe(ctx, userId, groupId)
	if err != nil {
		return errors.Wrap(err, "failed subscribe")
	}

	return nil
}

func (uc *UseCase) UnSubscribe(ctx context.Context, userId int64, groupId int64) error {
	ctx = uc.Logger.WithCaller(ctx)

	err := uc.profileRepository.UnSubscribe(ctx, userId, groupId)
	if err != nil {
		return errors.Wrap(err, "failed subscribe")
	}

	return nil
}
