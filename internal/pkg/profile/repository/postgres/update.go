package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r profileRepository) UpdateProfile(ctx context.Context, profile models.ChangeProfile) error {
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
