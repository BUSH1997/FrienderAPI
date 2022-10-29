package profile

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	UpdateProfile(ctx context.Context, profile models.ChangeProfile) error
	CheckUserExists(ctx context.Context, user int64) (bool, error)
	Create(ctx context.Context, user int64, isGroup bool) error
	Subscribe(ctx context.Context, userId int64, groupId int64) error
	UnSubscribe(ctx context.Context, userId int64, groupId int64) error
}
