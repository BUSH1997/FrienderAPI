package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (u eventUsecase) GetAllPublic(ctx context.Context) ([]models.Event, error) {
	events, err := u.Events.GetAllPublic(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all public events in usecase")
	}

	return events, nil
}

func (u eventUsecase) GetEventById(ctx context.Context, id string) (models.Event, error) {
	return u.Events.GetEventById(ctx, id)
}

func (u eventUsecase) GetUserEvents(ctx context.Context, id int64) ([]models.Event, error) {
	return u.Events.GetUserEvents(ctx, id)
}

func (u eventUsecase) GetAll(ctx context.Context, filter event.FilterGetAll) ([]models.Event, error) {
	events, err := u.Events.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all events in usecase")
	}

	return events, nil
}

func (u eventUsecase) GetAllCategories(ctx context.Context) ([]string, error) {
	return u.Events.GetAllCategories(ctx)
}
