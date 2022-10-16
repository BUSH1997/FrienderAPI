package event

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type FilterGetAll struct {
	IsSubscribe bool
	IsActive    bool
	User        int
	IsOwner     bool
	Page        int
	Limit       int
}

type Usecase interface {
	Create(ctx context.Context, event models.Event) (models.Event, error)
	Update(ctx context.Context, event models.Event) error
	Get(ctx context.Context, params models.GetEventParams) ([]models.Event, error)
	GetAllPublic(ctx context.Context) ([]models.Event, error)
	GetAll(ctx context.Context, filter FilterGetAll) ([]models.Event, error)
	GetEventById(ctx context.Context, id string) (models.Event, error)
	GetUserEvents(ctx context.Context, id int64) ([]models.Event, error)
	SubscribeEvent(ctx context.Context, event string) error
	UnsubscribeEvent(ctx context.Context, event string) error
	Delete(ctx context.Context, event string) error
	Change(ctx context.Context, event models.Event) error
	GetAllCategories(ctx context.Context) ([]string, error)
}
