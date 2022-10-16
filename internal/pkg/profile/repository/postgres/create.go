package postgres

import (
	"context"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r profileRepository) Create(ctx context.Context, user int64) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		dbUser := db_models.User{
			Uid:           int(user),
			CurrentStatus: 1,
		}
		res := r.db.Create(&dbUser)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to create user")
		}

		dbUnlockedStatus := db_models.UnlockedStatus{
			UserID:   int(dbUser.ID),
			StatusID: 1,
		}
		res = r.db.Create(&dbUnlockedStatus)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to create unlocked status for user")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}
	
	return nil
}
