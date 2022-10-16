package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (uc *UseCase) UpdateProfile(ctx context.Context, profile models.ChangeProfile) error {
	err := uc.profileRepository.UpdateProfile(ctx, profile)
	if err != nil {
		return errors.Wrap(err, "failed to update profile")
	}

	return nil
}

func (uc *UseCase) ChangeEventPriority(ctx context.Context, eventPriority models.UidEventPriority) error {
	err := uc.eventRepository.UpdateEventPriority(ctx, eventPriority)
	if err != nil {
		return errors.Wrap(err, "failed to update event priority")
	}

	return nil
}
