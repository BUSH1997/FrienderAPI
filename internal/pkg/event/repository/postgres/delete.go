package postgres

import (
	"context"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	event_pkg "github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/gorm"
)

func (r eventRepository) Delete(ctx context.Context, event string, groupInfo models.GroupInfo) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvent db_models.Event
		res := r.db.Take(&dbEvent, "uid = ?", event)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get event by uid %s", event)
		}

		userID := contextlib.GetUser(ctx)

		var dbUser db_models.User
		res = r.db.Take(&dbUser, "uid = ?", userID)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get user by uid %d", userID)
		}

		res = r.db.Model(&db_models.User{}).Where("id = ?", dbUser.ID).
			Update("created_events", dbUser.CreatedEventsCount-1)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to update user created events")
		}

		if groupInfo.GroupId == 0 {
			if dbEvent.Owner != int(dbUser.ID) {
				return event_pkg.ErrNoDeleteAccess.WithMessage("user is not events owner, cannot delete")
			}
		} else {
			//надо проверить является ли юзер админом
			var groups []db_models.Group
			res = r.db.Find(&groups, "user_id = ?", userID)
			if err := res.Error; err != nil {
				return errors.Wrap(err, "failed to get admin group")
			}

			checkAdmin := false
			for _, group := range groups {
				if int64(group.GroupId) == groupInfo.GroupId {
					checkAdmin = true
				}
			}

			if !checkAdmin {
				return event_pkg.ErrNoDeleteAccess.WithMessage("user is not admin, cannot delete")
			}
		}

		dbEventSharing := db_models.EventSharing{}
		res = r.db.Take(&dbEventSharing, "event_id = ? AND user_id = ?", dbEvent.ID, dbUser.ID)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get event sharing")
		}

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
