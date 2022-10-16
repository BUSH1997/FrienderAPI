package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r statusRepository) GetUserCurrentStatus(ctx context.Context, id int64) (models.Status, error) {
	var status models.Status

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbUser db_models.User
		res := r.db.Take(&dbUser, "uid = ?", id)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get user by uid %d", id)
		}

		var dbStatus db_models.Status
		res = r.db.Take(&dbStatus, "user_id = ?", dbUser.ID)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get status")
		}

		status.Title = dbStatus.Title
		status.Id = dbStatus.UID
		status.IsLocked = false

		return nil
	})
	if err != nil {
		return models.Status{}, errors.Wrap(err, "failed to make transaction")
	}

	return status, nil
}

func (r statusRepository) GetAllUserStatuses(ctx context.Context, id int64) ([]models.Status, error) {
	var statuses []models.Status

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbStatuses []db_models.Status

		res := r.db.Model(&db_models.Status{}).
			Joins("JOIN unlocked_statuses on unlocked_statuses.status_id = statuses.id").
			Joins("JOIN users on users.id = unlocked_statuses.user_id").
			Find(&dbStatuses, "users.uid = ?", id)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user statuses")
		}

		statuses = make([]models.Status, 0, len(dbStatuses))
		for _, dbStatus := range dbStatuses {
			status := models.Status{
				Title:    dbStatus.Title,
				Id:       dbStatus.UID,
				IsLocked: false,
			}

			statuses = append(statuses, status)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return statuses, nil
}
