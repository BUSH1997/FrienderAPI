package revindex

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/stammer"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

func (r eventRepository) Create(ctx context.Context, event models.Event) error {
	terms, err := stammer.GetStammers(stammer.FilterSkipList(strings.Split(event.Title, " "), r.skipList))
	if err != nil {
		return errors.Wrap(err, "failed to get stammers from title")
	}

	dbRevindexEvent := db_models.RevindexEvent{
		UID: event.Uid,
	}
	res := r.db.Create(&dbRevindexEvent)
	if err = res.Error; err != nil {
		return errors.Wrap(err, "failed to create revindex event")
	}

	for _, term := range terms {
		err := r.createOrUpdateRevindex(term, int64(dbRevindexEvent.ID))
		if err != nil {
			return errors.Wrapf(err, "failed to create or update revindex for term %s", term)
		}
	}

	err = r.events.Create(ctx, event)
	if err != nil {
		return errors.Wrap(err, "failed to create event")
	}

	return nil
}

func (r eventRepository) createOrUpdateRevindex(term string, ID int64) error {
	var dbRevindexWords []db_models.RevindexWord

	res := r.db.Model(db_models.RevindexWord{}).
		Where("word = ?", term).
		Find(&dbRevindexWords)

	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to get revindex words")
	}

	if len(dbRevindexWords) == 0 {
		dbRevindexWord := db_models.RevindexWord{
			Word:   term,
			Events: pq.Int64Array{ID},
		}

		res := r.db.Create(&dbRevindexWord)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to create revindex word")
		}

		return nil
	}

	err := r.updateRevindexWord(term, dbRevindexWords[0].Events, ID)
	if err != nil {
		return errors.Wrap(err, "failed to update revindex word")
	}

	return nil
}

func (r eventRepository) updateRevindexWord(term string, eventIDs []int64, ID int64) error {
	eventIDList := append(eventIDs, ID)
	sort.Slice(eventIDList, func(i, j int) bool {
		return eventIDList[i] < eventIDList[j]
	})

	res := r.db.Model(&db_models.RevindexWord{}).
		Where("word = ?", term).
		Update("events", pq.Int64Array(eventIDList))

	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to update revindex word in db")
	}

	return nil
}
