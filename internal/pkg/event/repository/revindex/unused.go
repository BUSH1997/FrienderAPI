package revindex

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

func (r eventRepository) GetAllPublic(ctx context.Context) ([]models.Event, error) {
	return r.events.GetAllPublic(ctx)
}

func (r eventRepository) GetAll(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	return r.events.GetAll(ctx, params)
}

func (r eventRepository) GetOwnerEvents(ctx context.Context, user int64) ([]models.Event, error) {
	return r.events.GetOwnerEvents(ctx, user)
}

func (r eventRepository) GetEventById(ctx context.Context, id string) (models.Event, error) {
	return r.events.GetEventById(ctx, id)
}

func (r eventRepository) GetUserEvents(ctx context.Context, user int64) ([]models.Event, error) {
	return r.events.GetUserEvents(ctx, user)
}

func (r eventRepository) UploadImage(ctx context.Context, uid string, link string) error {
	return r.events.UploadImage(ctx, uid, link)
}

func (r eventRepository) GetAllCategories(ctx context.Context) ([]string, error) {
	return r.events.GetAllCategories(ctx)
}

func (r eventRepository) GetUserActiveEvents(ctx context.Context, user int64, params models.GetEventParams) ([]models.Event, error) {
	return r.events.GetUserActiveEvents(ctx, user, params)
}

func (r eventRepository) GetUserVisitedEvents(ctx context.Context, user int64) ([]models.Event, error) {
	return r.events.GetUserVisitedEvents(ctx, user)
}

func (r eventRepository) UploadAvatar(ctx context.Context, uid string, link string, vkId string) error {
	return r.events.UploadAvatar(ctx, uid, link, vkId)
}

func (r eventRepository) GetSubscriptionEvents(ctx context.Context, user int64) ([]models.Event, error) {
	return r.events.GetSubscriptionEvents(ctx, user)
}

func (r eventRepository) GetGroupEvent(ctx context.Context, group int64, isActive models.Bool, params models.GetEventParams) ([]models.Event, error) {
	return r.events.GetGroupEvent(ctx, group, isActive, params)
}

func (r eventRepository) GetGroupAdminEvent(
	ctx context.Context,
	group int64,
	isAdmin models.Bool,
	isActive models.Bool,
) ([]models.Event, error) {
	return r.events.GetGroupAdminEvent(ctx, group, isAdmin, isActive)
}

func (r eventRepository) UpdateEventPriority(ctx context.Context, eventPriority models.UidEventPriority) error {
	return r.events.UpdateEventPriority(ctx, eventPriority)
}

func (r eventRepository) Subscribe(ctx context.Context, event string) error {
	return r.events.Subscribe(ctx, event)
}

func (r eventRepository) UnSubscribe(ctx context.Context, event string) error {
	return r.events.UnSubscribe(ctx, event)
}
