package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r eventRepository) GetAllPublic(ctx context.Context) ([]models.Event, error) {
	var dbEvents []Event

	res := r.db.Find(&dbEvents, "is_public = ?", true)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err := res.Error; err != nil {
		return nil, errors.Wrap(err, "failed to check user")
	}

	events := make([]models.Event, 0, len(dbEvents))
	for _, dbEvent := range dbEvents {
		event := models.Event{
			Uid:      dbEvent.Uid,
			Title:    dbEvent.Title,
			StartsAt: dbEvent.StartsAt,
			IsPublic: dbEvent.IsPublic,
		}

		events = append(events, event)
	}

	return events, nil
}

func (r eventRepository) GetEventById(ctx context.Context, id string) (models.Event, error) {
	return models.Event{}, nil
}
