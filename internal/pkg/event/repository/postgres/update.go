package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r eventRepository) Update(ctx context.Context, event models.Event) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		dbEvent := db_models.Event{
			Uid:      event.Uid,
			Title:    event.Title,
			StartsAt: event.StartsAt,
			IsPublic: event.IsPublic,
		}

		res := r.db.Model(&db_models.Event{}).Where("uid = ?", dbEvent.Uid).Updates(map[string]interface{}{
			"uid":       dbEvent.Uid,
			"title":     dbEvent.Title,
			"starts_at": dbEvent.StartsAt,
			"is_public": dbEvent.IsPublic,
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
