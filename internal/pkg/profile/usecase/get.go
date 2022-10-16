package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (uc *UseCase) GetOneProfile(ctx context.Context, id int64) (models.Profile, error) {
	currentStatus, err := uc.statusRepository.GetUserCurrentStatus(ctx, id)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get profile status")
	}

	activeEvents, err := uc.eventRepository.GetUserActiveEvents(ctx, id)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get user active events")
	}

	visitedEvents, err := uc.eventRepository.GetUserActiveEvents(ctx, id)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get user visited events")
	}

	awards, err := uc.awardRepository.GetUserAwards(ctx, id)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get user awards")
	}

	profile := models.Profile{
		ProfileStatus: currentStatus,
		Awards:        awards,
		ActiveEvents:  activeEvents,
		VisitedEvents: visitedEvents,
	}

	return profile, nil
}

func (uc *UseCase) GetAllProfileStatuses(ctx context.Context, id int64) ([]models.Status, error) {
	statuses, err := uc.statusRepository.GetAllUserStatuses(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all profile statuses")
	}

	return statuses, nil
}
