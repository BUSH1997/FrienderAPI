package event

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	Create(ctx context.Context, event models.Event) error
	Update(ctx context.Context, event models.Event) error
	GetAllPublic(ctx context.Context) ([]models.Event, error)
	GetAll(ctx context.Context, params models.GetEventParams) ([]models.Event, error)
	GetOwnerEvents(ctx context.Context, user int64) ([]models.Event, error)
	GetEventById(ctx context.Context, id string) (models.Event, error)
	GetUserEvents(ctx context.Context, user int64) ([]models.Event, error)
	UploadImage(ctx context.Context, uid string, link string) error
	UploadAvatar(ctx context.Context, uid string, link string, vkId string) error
	GetAllCategories(ctx context.Context) ([]string, error)
	GetUserActiveEvents(ctx context.Context, user int64) ([]models.Event, error)
	GetUserVisitedEvents(ctx context.Context, user int64) ([]models.Event, error)
	GetSubscriptionEvents(ctx context.Context, user int64) ([]models.Event, error)
	GetGroupEvent(ctx context.Context, group int64, isActive models.Bool) ([]models.Event, error)
	GetGroupAdminEvent(ctx context.Context, group int64, isAdmin models.Bool, isActive models.Bool) ([]models.Event, error)
	UpdateEventPriority(ctx context.Context, eventPriority models.UidEventPriority) error
	Subscribe(ctx context.Context, event string) error
	UnSubscribe(ctx context.Context, event string) error
	Delete(ctx context.Context, event string, groupInfo models.GroupInfo) error
}
