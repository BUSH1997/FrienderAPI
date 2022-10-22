package group

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	Create(ctx context.Context, group models.Group) error
	GetAdministeredGroupByUserId(ctx context.Context, userId int) ([]models.Group, error)
	CheckIfAdmin(ctx context.Context, userId int, groupId int) (bool, error)
}
