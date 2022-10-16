package profile

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type UseCase interface {
	GetOneProfile(ctx context.Context, id int) (models.Profile, error)
	GetAllProfileStatuses(ctx context.Context, id int) ([]models.Status, error)
	UpdateProfile(ctx context.Context, profile models.ChangeProfile) error
	ChangeEventPriority(ctx context.Context, eventPriority models.UidEventPriority) error
}
