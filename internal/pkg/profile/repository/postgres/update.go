package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/gorm"
)

func (r profileRepository) UpdateProfile(ctx context.Context, profile models.ChangeProfile) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		res := r.db.Model(&db_models.User{}).Where("uid = ?", profile.ProfileId).
			Updates(map[string]interface{}{
				"current_status": profile.Status.Id,
			})
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update user %d", profile.ProfileId)
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func (r profileRepository) Subscribe(ctx context.Context, userId int64, groupId int64) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbUser db_models.User
		res := r.db.Take(&dbUser, "uid = ?", userId)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user by id")
		}

		var dbProfile db_models.User
		res = r.db.Take(&dbProfile, "uid = ?", groupId)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user by id")
		}

		dbSubscribe := db_models.SubscribeProfileSharing{
			ProfileId: int64(dbProfile.ID),
			UserId:    int64(dbUser.ID),
		}
		res = r.db.Create(&dbSubscribe)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to create subscribe")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func (r profileRepository) UnSubscribe(ctx context.Context, userId int64, groupId int64) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbUser db_models.User
		res := r.db.Take(&dbUser, "uid = ?", userId)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user by id")
		}

		var dbProfile db_models.User
		res = r.db.Take(&dbProfile, "uid = ?", groupId)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user by id")
		}

		dbSubscribe := db_models.SubscribeProfileSharing{
			ProfileId: int64(dbProfile.ID),
			UserId:    int64(dbUser.ID),
		}
		res = r.db.Delete(&dbSubscribe, "profile_id = ? and user_id = ?", dbSubscribe.ProfileId, dbSubscribe.UserId)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to delete subscribe")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}
