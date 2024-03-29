package revindex

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

func (r eventRepository) GetAll(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	return r.events.GetAll(ctx, params)
}

func (r eventRepository) CheckIfExists(ctx context.Context, event models.Event) (bool, error) {
	return r.events.CheckIfExists(ctx, event)
}

func (r eventRepository) GetCountEvents(ctx context.Context, typeEvents string) (int64, error) {
	return r.events.GetCountEvents(ctx, typeEvents)
}

func (r eventRepository) GetEventById(ctx context.Context, id string) (models.Event, error) {
	return r.events.GetEventById(ctx, id)
}

func (r eventRepository) UploadImage(ctx context.Context, uid string, link string) error {
	return r.events.UploadImage(ctx, uid, link)
}

func (r eventRepository) GetAllCategories(ctx context.Context) ([]string, error) {
	return r.events.GetAllCategories(ctx)
}

func (r eventRepository) UploadAvatar(ctx context.Context, uid string, link string, vkId string) error {
	return r.events.UploadAvatar(ctx, uid, link, vkId)
}

func (r eventRepository) GetSubscriptionEvents(ctx context.Context, user int64) ([]models.Event, error) {
	return r.events.GetSubscriptionEvents(ctx, user)
}

func (r eventRepository) GetGroupEvents(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	return r.events.GetGroupEvents(ctx, params)
}

func (r eventRepository) GetSharings(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	return r.events.GetSharings(ctx, params)
}

func (r eventRepository) UpdateEventPriority(ctx context.Context, eventPriority models.UidEventPriority) error {
	return r.events.UpdateEventPriority(ctx, eventPriority)
}

func (r eventRepository) Subscribe(ctx context.Context, event string) error {
	return r.events.Subscribe(ctx, event)
}

func (r eventRepository) UnSubscribe(ctx context.Context, event string, user int64) error {
	return r.events.UnSubscribe(ctx, event, user)
}

func (r eventRepository) AddAlbum(ctx context.Context, eventUid string, albumUid string) error {
	return r.events.AddAlbum(ctx, eventUid, albumUid)
}

func (r eventRepository) DeleteAlbum(ctx context.Context, eventUid string, albumUid string) error {
	return r.events.DeleteAlbum(ctx, eventUid, albumUid)
}
