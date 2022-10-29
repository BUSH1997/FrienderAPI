package postgres

import (
	"context"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r searchRepository) GetEventUIDs(ctx context.Context, terms []string) ([]string, error) {
	var eventIDsArray [][]int64
	for _, term := range terms {
		var dbRevindexWord db_models.RevindexWord
		res := r.db.Model(&db_models.RevindexWord{}).
			Where("word = ?", term).
			Take(&dbRevindexWord)
		if err := res.Error; err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		if err := res.Error; err != nil {
			return nil, errors.Wrapf(err, "failed to get revindex for word %s", term)
		}

		eventIDsArray = append(eventIDsArray, dbRevindexWord.Events)
	}

	eventIDsMap := make(map[int64]int, 0)
	for _, eventIDs := range eventIDsArray {
		for _, eventID := range eventIDs {
			eventIDsMap[eventID] += 1
		}
	}

	var eventIDs []int64
	for k, v := range eventIDsMap {
		if v == len(terms) {
			eventIDs = append(eventIDs, k)
		}
	}

	if len(eventIDs) == 0 {
		return nil, nil
	}

	var dbRevindexEvents []db_models.RevindexEvent
	res := r.db.Model(&db_models.RevindexEvent{}).
		Find(&dbRevindexEvents, eventIDs)
	if err := res.Error; err != nil {
		return nil, errors.Wrap(err, "failed to get revindex events")
	}

	eventUIDs := make([]string, 0, len(dbRevindexEvents))
	for _, dbRevindexEvent := range dbRevindexEvents {
		eventUIDs = append(eventUIDs, dbRevindexEvent.UID)
	}

	return eventUIDs, nil
}
