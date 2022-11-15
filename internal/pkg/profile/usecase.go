package profile

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type UseCase interface {
	GetOneProfile(ctx context.Context, id int64) (models.Profile, error)
	GetAllProfileStatuses(ctx context.Context) ([]models.Status, error)
	UpdateProfile(ctx context.Context, profile models.ChangeProfile) error
	ChangeEventPriority(ctx context.Context, eventPriority models.UidEventPriority) error
	Subscribe(ctx context.Context, userId int64, groupId int64) error
	UnSubscribe(ctx context.Context, userId int64, groupId int64) error
	GetSubscribe(cxt context.Context, userId int64) (models.Subscriptions, error)
	GetFriends(ctx context.Context, userId int64) ([]int64, error)
	GetCities(ctx context.Context) ([]string, error)
}
