package group

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type UseCase interface {
	Create(ctx context.Context, group models.GroupInput) error
	Update(ctx context.Context, group models.GroupInput) error
	GetAdministeredGroupByUserId(ctx context.Context, userId string) ([]models.Group, error)
	Get(ctx context.Context, userId int64) (models.Group, error)
	CheckIfAdmin(ctx context.Context, userId string, groupId int64) (bool, error)
}
