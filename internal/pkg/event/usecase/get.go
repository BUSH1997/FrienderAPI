package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (uc eventUsecase) GetAllPublic(ctx context.Context) ([]models.Event, error) {
	events, err := uc.Events.GetAllPublic(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all public events in usecase")
	}

	return events, nil
}

func (uc eventUsecase) GetEventById(ctx context.Context, id string) (models.Event, error) {
	return uc.Events.GetEventById(ctx, id)
}

func (uc eventUsecase) GetUserEvents(ctx context.Context, id int64) ([]models.Event, error) {
	return uc.Events.GetUserEvents(ctx, id)
}

func (uc eventUsecase) GetAll(ctx context.Context, filter event.FilterGetAll) ([]models.Event, error) {
	events, err := uc.Events.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all events in usecase")
	}

	return events, nil
}

func (uc eventUsecase) GetAllCategories(ctx context.Context) ([]string, error) {
	return uc.Events.GetAllCategories(ctx)
}

func (uc eventUsecase) Get(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	if params.IsOwner.IsDefinedTrue() {
		return uc.Events.GetOwnerEvents(ctx, params.UserID)
	}
	if params.IsActive.IsDefinedTrue() {
		return uc.Events.GetUserActiveEvents(ctx, params.UserID)
	}
	if params.IsActive.IsDefinedFalse() {
		return uc.Events.GetUserVisitedEvents(ctx, params.UserID)
	}
	if params.IsSubscriber.IsDefinedTrue() {
		return uc.Events.GetSubscriptionEvents(ctx, params.UserID)
	}

	return uc.Events.GetAll(ctx)
}
