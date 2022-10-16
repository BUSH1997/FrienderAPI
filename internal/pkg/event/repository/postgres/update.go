package postgres

import (
	"context"
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func (r eventRepository) Update(ctx context.Context, event models.Event) error {
	fmt.Println(event.Uid)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		res := r.db.Model(&db_models.Event{}).Where("uid = ?", event.Uid).Updates(map[string]interface{}{
			"uid":          event.Uid,
			"title":        event.Title,
			"description":  event.Description,
			"starts_at":    event.StartsAt,
			"time_updated": time.Now().Unix(),
			"geo":          strconv.Itoa(int(event.GeoData.Longitude)) + "," + strconv.Itoa(int(event.GeoData.Latitude)),
			"is_public":    event.IsPublic,
			"is_private":   event.IsPrivate,
		})
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event, uid %d", event.Uid)
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func (r eventRepository) UpdateEventPriority(ctx context.Context, eventPriority models.UidEventPriority) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		oldPriority, err := r.getPriority(ctx, eventPriority.UidUser, eventPriority.UidEvent)
		if err != nil {
			return errors.Wrap(err, "failed to get old event priority")
		}

		var dbUser db_models.User
		res := r.db.Take(&dbUser, "uid = ?", eventPriority.UidUser)
		if err = res.Error; err != nil {
			return errors.Wrap(err, "failed to get user")
		}

		posDiff := 0
		if oldPriority-eventPriority.Priority > 0 {
			posDiff = 1
		} else {
			posDiff = -1
		}

		res = r.db.Model(&db_models.EventSharing{}).
			Where("user_id = ? AND priority BETWEEN ? AND ?", dbUser.ID, oldPriority-posDiff, eventPriority.Priority).
			Update("priority", gorm.Expr("event_sharings.priority + ?", posDiff))
		if err = res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharings priority")
		}

		res = r.db.Model(&db_models.EventSharing{}).
			Where("user_id = ? AND priority = ?", dbUser.ID, oldPriority).
			Update("priority", eventPriority.Priority)
		if err = res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharing priority")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func (r eventRepository) Subscribe(ctx context.Context, user int64, event string) error {
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
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			dbEventSharing.EventID = int(dbEvent.ID)
			dbEventSharing.UserID = int(dbUser.ID)
			res = r.db.Create(&dbEventSharing)
			if err := res.Error; err != nil {
				return errors.Wrapf(err, "failed to create event sharing")
			}
		}
		if err := res.Error; err != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "failed to get event sharing")
		}

		var dbEventSharings []db_models.EventSharing
		res = r.db.Order("priority desc").Find(&dbEventSharings, "user_id = ?", dbUser.ID)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get event sharings")
		}

		priority := 1
		if len(dbEventSharings) > 0 {
			priority = dbEventSharings[0].Priority
		}

		res = r.db.Model(&db_models.EventSharing{}).
			Where("event_id = ? AND user_id = ?", dbEvent.ID, dbUser.ID).
			Updates(map[string]interface{}{
				"is_deleted": false,
				"priority":   priority,
			})
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharing")
		}

		res = r.db.Model(&db_models.Event{}).
			Where("id = ?", dbEvent.ID).
			Update("count_members", gorm.Expr("events.count_members + ?", 1))
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event members count")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func (r eventRepository) UnSubscribe(ctx context.Context, user int64, event string) error {
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
			Where("user_id = ? AND priority > ?", currentPriority).
			Update("priority", gorm.Expr("event_sharings.priority + ?", -1))
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharings priority")
		}

		res = r.db.Model(&db_models.EventSharing{}).
			Where("event_id = ? AND user_id = ?", dbEvent.ID, dbUser.ID).
			Update("is_deleted", true)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharing")
		}

		res = r.db.Model(&db_models.Event{}).
			Where("id = ?", dbEvent.ID).
			Update("count_members", gorm.Expr("events.count_members + ?", -1))
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event members count")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

//res := r.db.Model(&db_models.EventSharing{}).
//	Joins("JOIN events on event_sharings.event_id = events.id").
//	Joins("JOIN users on event_sharings.user_id = users.id").
//	Where("users.uid = ? AND event.uid = ?", user, event).
//	Update("is", eventPriority.Priority)
//if err = res.Error; err != nil {
//	return errors.Wrapf(err, "failed to update event sharing priority")
//}
