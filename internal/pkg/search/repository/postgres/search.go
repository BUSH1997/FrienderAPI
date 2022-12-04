package postgres

import (
	"context"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
)

func (r searchRepository) GetEventUIDs(ctx context.Context, terms []string) ([]string, error) {
	ctx = r.logger.WithCaller(ctx)

	var eventIDsArray []string
	for _, term := range terms {
		var dbRevindexWord db_models.RevindexWord
		res := r.db.Model(&db_models.RevindexWord{}).
			Where("word LIKE ?", "%"+term+"%").
			Take(&dbRevindexWord)
		if err := res.Error; err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		if err := res.Error; err != nil {
			return nil, errors.Wrapf(err, "failed to get revindex for word %s", term)
		}

		eventIDsArray = append(eventIDsArray, dbRevindexWord.Events...)
	}

	eventIDsMap := make(map[int64]int, 0)
	eventWordPositionMap := make(map[int64]int64, 0)
	for _, searchEvent := range eventIDsArray {
		eventIDString := strings.Split(searchEvent, ":")[0]
		eventID, err := strconv.ParseInt(eventIDString, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse event id %s", eventIDString)
		}
		wordPositionString := strings.Split(searchEvent, ":")[1]
		wordPosition, err := strconv.ParseInt(wordPositionString, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse word position %s", wordPositionString)
		}

		eventIDsMap[eventID] += 1
		eventWordPositionMap[eventID] += wordPosition
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

	sort.Slice(dbRevindexEvents, func(i, j int) bool {
		return eventWordPositionMap[int64(dbRevindexEvents[i].ID)] < eventWordPositionMap[int64(dbRevindexEvents[j].ID)]
	})

	eventUIDs := make([]string, 0, len(dbRevindexEvents))
	for _, dbRevindexEvent := range dbRevindexEvents {
		eventUIDs = append(eventUIDs, dbRevindexEvent.UID)
	}

	return eventUIDs, nil
}
