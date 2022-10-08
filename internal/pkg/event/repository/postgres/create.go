package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (r eventRepository) Create(ctx context.Context, event models.Event) error {
	dbEvent := Event{
		Uid:      event.Uid,
		Title:    event.Title,
		StartsAt: event.StartsAt,
		IsPublic: event.IsPublic,
	}

	res := r.db.Create(&dbEvent)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to create event, uid %d", event.Uid)
	}

	return nil
}
