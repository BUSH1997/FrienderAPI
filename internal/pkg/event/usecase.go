package event

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Usecase interface {
	Create(ctx context.Context, event models.Event) error
	Update(ctx context.Context, event models.Event) error
	GetAllPublic(ctx context.Context) ([]models.Event, error)
	GetAll(ctx context.Context) ([]models.Event, error)
	GetEventById(ctx context.Context, id string) (models.Event, error)
	GetUserEvents(ctx context.Context, id string) ([]models.Event, error)
}