package postgres

import (
	"context"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r eventRepository) Delete(ctx context.Context, user int64, event string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvent db_models.Event
		res := r.db.Take(&dbEvent, "uid = ?", event)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get event by uid %s", event)
		}

		var dbUser db_models.User
		res = r.db.Take(&dbUser, "uid = ?", user)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get user by uid %d", user)
		}

		dbEventSharing := db_models.EventSharing{}
		res = r.db.Take(&dbEventSharing, "event_id = ? AND user_id = ?", dbEvent.ID, dbUser.ID)

		currentPriority := dbEventSharing.Priority

		res = r.db.Model(&db_models.EventSharing{}).
			Where("user_id = ? AND priority > ?", dbUser.ID, currentPriority).
			Update("priority", gorm.Expr("event_sharings.priority + ?", -1))
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharings priority")
		}

		res = r.db.Model(&db_models.EventSharing{}).
			Where("event_id = ?", dbEvent.ID).
			Update("is_deleted", true)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharings")
		}

		res = r.db.Model(&db_models.Event{}).
			Where("id = ?", dbEvent.ID).
			Update("is_deleted", true)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}
