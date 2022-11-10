package postgres

import (
	"context"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
)

func (r chatRepository) UpdateLastCheckTime(ctx context.Context, event string, user int64, time int64) error {
	var dbEvent db_models.Event
	res := r.db.Take(&dbEvent, "uid = ?", event)
	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to get event by id")
	}

	var dbUser db_models.User
	res = r.db.Take(&dbUser, "uid = ?", user)
	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to get user by uid")
	}

	res = r.db.Model(&db_models.EventSharing{}).
		Where("event_id = ?", dbEvent.ID).
		Where("user_id = ?", dbUser.ID).
		Update("time_last_check", time)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to update last check time for chat %s", event)
	}

	return nil
}
