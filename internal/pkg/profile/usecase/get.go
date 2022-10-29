package usecase

import (
	"context"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (uc *UseCase) GetOneProfile(ctx context.Context, userID int64) (models.Profile, error) {
	currentStatus, err := uc.statusRepository.GetUserCurrentStatus(ctx, userID)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get profile status")
	}

	activeEvents, err := uc.eventRepository.GetUserActiveEvents(ctx, userID)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get user active events")
	}

	visitedEvents, err := uc.eventRepository.GetUserActiveEvents(ctx, userID)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get user visited events")
	}

	awards, err := uc.awardRepository.GetUserAwards(ctx, userID)
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

func (uc *UseCase) GetAllProfileStatuses(ctx context.Context) ([]models.Status, error) {
	userID := contextlib.GetUser(ctx)

	statuses, err := uc.statusRepository.GetAllUserStatuses(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all profile statuses")
	}

	return statuses, nil
}

func (uc *UseCase) GetSubscribe(cxt context.Context, userId int64) ([]int, error) {
	subscribe, err := uc.profileRepository.GetSubscribe(cxt, userId)
	if err != nil {
		uc.Logger.WithError(err).Errorf("[GetSubscribe] failed get subscribe")
		return []int{}, nil
	}

	result := make([]int, len(subscribe))
	for _, v := range subscribe {
		result = append(result, int(v.Id))
	}

	return result, nil
}
