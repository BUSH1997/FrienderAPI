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
	GetSubscribe(cxt context.Context, userId int64) ([]models.SubscribeType, error)
	GetCities(ctx context.Context) ([]string, error)
	GetOneProfile(ctx context.Context, userID int64) (models.Profile, error)
	GetAllUserStatuses(ctx context.Context, id int64) ([]models.Status, error)
}
