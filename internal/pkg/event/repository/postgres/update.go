package postgres

import (
	"context"
	"fmt"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	event_pkg "github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/gorm"
	"math"
	"time"
)

func (r eventRepository) Update(ctx context.Context, event models.Event) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		res := r.db.Model(&db_models.Event{}).Where("uid = ?", event.Uid).Updates(map[string]interface{}{
			"uid":           event.Uid,
			"title":         event.Title,
			"description":   event.Description,
			"starts_at":     event.StartsAt,
			"time_updated":  time.Now().Unix(),
			"geo":           fmt.Sprintf("%f;;%f;;%s", event.GeoData.Longitude, event.GeoData.Latitude, event.GeoData.Address),
			"is_public":     event.IsPublic,
			"is_private":    event.IsPrivate,
			"ticket":        fmt.Sprintf("%s;;%s", event.Ticket.Link, event.Ticket.Cost),
			"members_limit": event.MembersLimit,
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
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEventSharing db_models.EventSharing

		res := r.db.Model(&db_models.EventSharing{}).
			Joins("JOIN users on event_sharings.user_id = users.id").
			Joins("JOIN events on event_sharings.event_id = events.id").
			Take(&dbEventSharing, "users.uid = ? AND events.uid = ?",
				eventPriority.UidUser, eventPriority.UidEvent)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get event sharing")
		}

		oldPriority := dbEventSharing.Priority

		var dbUser db_models.User
		res = r.db.Take(&dbUser, "uid = ?", eventPriority.UidUser)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user")
		}

		smallerPriority := int(math.Min(float64(oldPriority), float64(eventPriority.Priority)))
		greaterPriority := int(math.Max(float64(oldPriority), float64(eventPriority.Priority)))

		posDiff := 0
		if oldPriority-eventPriority.Priority > 0 {
			posDiff = 1
		} else {
			posDiff = -1
		}

		res = r.db.Model(&db_models.EventSharing{}).
			Where("user_id = ?", dbUser.ID).
			Where("priority BETWEEN ? AND ?", smallerPriority, greaterPriority).
			Update("priority", gorm.Expr("event_sharings.priority + ?", posDiff))
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharings priority")
		}

		res = r.db.Model(&db_models.EventSharing{}).
			Where("id = ?", dbEventSharing.ID).
			Update("priority", eventPriority.Priority)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharing priority")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func (r eventRepository) Subscribe(ctx context.Context, event string) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvent db_models.Event
		res := r.db.Take(&dbEvent, "uid = ?", event)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get event by uid %s", event)
		}

		userID := contextlib.GetUser(ctx)
		for _, banned := range dbEvent.BlackList {
			if banned == userID {
				return event_pkg.ErrNoAccessForBanned
			}
		}

		var dbUser db_models.User
		res = r.db.Take(&dbUser, "uid = ?", userID)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get user by uid %d", userID)
		}

		var sharingsExist []db_models.EventSharing
		res = r.db.Model(&db_models.EventSharing{}).
			Find(&sharingsExist, "user_id = ? AND is_deleted = ?", dbUser.ID, false)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user event sharings")
		}

		currentMaxPriority := len(sharingsExist)

		dbEventSharing := db_models.EventSharing{}
		res = r.db.Take(&dbEventSharing, "event_id = ? AND user_id = ?", dbEvent.ID, dbUser.ID)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			dbEventSharing.EventID = int(dbEvent.ID)
			dbEventSharing.UserID = int(dbUser.ID)
			res = r.db.Create(&dbEventSharing)
			if err := res.Error; err != nil {
				if postgres.ProcessError(err) == postgres.UniqueViolationError {
					err = errors.Transform(err, event_pkg.ErrAlreadyExists)
				}

				return errors.Wrapf(err, "failed to create event sharing")
			}
		}
		if err := res.Error; err != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "failed to get event sharing")
		}

		res = r.db.Model(&db_models.EventSharing{}).
			Where("event_id = ? AND user_id = ?", dbEvent.ID, dbUser.ID).
			Updates(map[string]interface{}{
				"is_deleted": false,
				"priority":   currentMaxPriority + 1,
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

func (r eventRepository) UnSubscribe(ctx context.Context, event string, user int64) error {
	ctx = r.logger.WithCaller(ctx)

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

		var dbOwner db_models.User
		res = r.db.Take(&dbOwner, "id = ?", dbEvent.Owner)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get owner")
		}

		if user != contextlib.GetUser(ctx) && user != int64(dbOwner.Uid) {
			res = r.db.Model(&db_models.Event{}).
				Where("id = ?", dbEvent.ID).
				Update("blacklist", gorm.Expr("array_append(events.blacklist, ?)", user))
			if err := res.Error; err != nil {
				return errors.Wrapf(err, "failed to update event blacklist")
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func (r eventRepository) AddAlbum(ctx context.Context, eventUid string, albumUid string) error {
	ctx = r.logger.WithCaller(ctx)

	var dbEvent db_models.Event

	res := r.db.Take(&dbEvent, "uid = ?", eventUid)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to get event by uid %s", eventUid)
	}

	dbEvent.Albums = append(dbEvent.Albums, albumUid)
	res = r.db.Model(&db_models.Event{}).
		Where("uid = ?", eventUid).
		Update("albums", dbEvent.Albums)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "error update albums event")
	}

	return nil
}

func (r eventRepository) DeleteAlbum(ctx context.Context, eventUid string, albumUid string) error {
	ctx = r.logger.WithCaller(ctx)

	var dbEvent db_models.Event

	res := r.db.Take(&dbEvent, "uid = ?", eventUid)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to get event by uid %s", eventUid)
	}

	for i, v := range dbEvent.Albums {
		if v == albumUid {
			dbEvent.Albums = append(dbEvent.Albums[:i], dbEvent.Albums[i+1:]...)
			break
		}
	}
	res = r.db.Model(&db_models.Event{}).
		Where("uid = ?", eventUid).
		Update("albums", dbEvent.Albums)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "error update albums event")
	}

	return nil
}
