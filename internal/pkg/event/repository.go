package event

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	Create(ctx context.Context, event models.Event) error
	Update(ctx context.Context, event models.Event) error
	GetAll(ctx context.Context, params models.GetEventParams) ([]models.Event, error)
	GetEventById(ctx context.Context, id string) (models.Event, error)
	UploadImage(ctx context.Context, uid string, link string) error
	UploadAvatar(ctx context.Context, uid string, link string, vkId string) error
	GetAllCategories(ctx context.Context) ([]string, error)
	GetSharings(ctx context.Context, params models.GetEventParams) ([]models.Event, error)
	GetSubscriptionEvents(ctx context.Context, user int64) ([]models.Event, error)
	GetGroupEvents(ctx context.Context, params models.GetEventParams) ([]models.Event, error)
	UpdateEventPriority(ctx context.Context, eventPriority models.UidEventPriority) error
	Subscribe(ctx context.Context, event string) error
	UnSubscribe(ctx context.Context, event string) error
	Delete(ctx context.Context, event string, groupInfo models.GroupInfo) error
	AddAlbum(ctx context.Context, eventUid string, albumUid string) error
	DeleteAlbum(ctx context.Context, eventUid string, albumUid string) error
}
