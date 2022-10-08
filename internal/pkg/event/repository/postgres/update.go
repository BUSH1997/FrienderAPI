package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
)

func (r eventRepository) Update(ctx context.Context, event models.Event) error {
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
}
