package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r profileRepository) CheckUserExists(ctx context.Context, user int64) (bool, error) {
	userExists := true
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbUser db_models.User
		res := r.db.Take(&dbUser, "uid = ?", user)
		if err := res.Error; errors.Is(err, gorm.ErrRecordNotFound) {
			userExists = false
			return nil
		}
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user by id")
		}

		return nil
	})
	if err != nil {
		return false, errors.Wrap(err, "failed to make transaction")
	}

	return userExists, nil
}

func (r profileRepository) GetSubscribe(cxt context.Context, userId int64) ([]models.SubscribeType, error) {
	var result []models.SubscribeType
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var currentUser db_models.User
		res := r.db.Take(&currentUser, "uid = ?", userId)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get subscribe")
		}

		var dbSubscribers []db_models.SubscribeProfileSharing

		res = r.db.Find(&dbSubscribers, "user_id = ?", currentUser.ID)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get subscribe")
		}

		for _, dbSubscribe := range dbSubscribers {
			var user db_models.User
			res := r.db.Take(&user, "id = ?", dbSubscribe.ProfileId)
			if err := res.Error; err != nil {
				return errors.Wrap(err, "failed to get subscribe")
			}

			subscribeType := models.SubscribeType{
				Id:      dbSubscribe.ProfileId,
				IsGroup: user.IsGroup,
			}

			result = append(result, subscribeType)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return result, nil

}
