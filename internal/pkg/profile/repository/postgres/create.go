package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/gorm"
)

func (r profileRepository) Create(ctx context.Context, user int64, isGroup bool) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		dbUser := db_models.User{
			Uid:           int(user),
			CurrentStatus: 1,
			IsGroup:       isGroup,
		}
		res := r.db.Create(&dbUser)
		if err := res.Error; err != nil {
			if postgres.ProcessError(err) == postgres.UniqueViolationError {
				err = errors.Transform(err, profile.ErrAlreadyExists)
			}

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

		dbUAuth := db_models.AuthUser{
			UID: user,
		}

		res = r.db.Create(&dbUAuth)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to create auth user")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}
