package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
	"sort"
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

func (uc eventUsecase) GetAllCategories(ctx context.Context) ([]string, error) {
	return uc.Events.GetAllCategories(ctx)
}

func (uc eventUsecase) Get(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	events, err := uc.routerGet(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get events in usecase")
	}

	sort.SliceStable(events, func(i, j int) bool {
		return events[i].StartsAt > events[j].StartsAt
	})

	if events == nil {
		events = make([]models.Event, 0)
	}

	return events, nil
}

func (uc eventUsecase) routerGet(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	if params.IsOwner.IsDefinedTrue() {
		return uc.Events.GetOwnerEvents(ctx, params.UserID)
	}
	if params.IsActive.IsDefinedTrue() && params.UserID != 0 {
		return uc.Events.GetUserActiveEvents(ctx, params.UserID)
	}
	if params.IsActive.IsDefinedFalse() && params.UserID != 0 {
		return uc.Events.GetUserVisitedEvents(ctx, params.UserID)
	}
	if params.IsSubscriber.IsDefinedTrue() {
		return uc.Events.GetSubscriptionEvents(ctx, params.UserID)
	}
	if params.GroupId != 0 && (params.IsAdmin.IsDefinedTrue() || params.IsAdmin.IsDefinedFalse()) {
		return uc.Events.GetGroupAdminEvent(ctx, params.GroupId, params.IsAdmin, params.IsActive)
	}
	if params.GroupId != 0 {
		return uc.Events.GetGroupEvent(ctx, params.GroupId, params.IsActive)
	}
	return uc.Events.GetAll(ctx, params)
}
