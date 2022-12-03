package revindex

import (
	"context"
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/stammer"
	"github.com/lib/pq"
	"sort"
	"strconv"
	"strings"
)

func (r eventRepository) Create(ctx context.Context, event models.Event) error {
	ctx = r.logger.WithCaller(ctx)

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

	for i, term := range terms {
		err := r.createOrUpdateRevindex(term, int64(dbRevindexEvent.ID), i)
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

func (r eventRepository) createOrUpdateRevindex(term string, ID int64, termPosition int) error {
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
			Events: pq.StringArray{fmt.Sprintf("%d:%d", ID, termPosition)},
		}

		res := r.db.Create(&dbRevindexWord)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to create revindex word")
		}

		return nil
	}

	err := r.updateRevindexWord(term, dbRevindexWords[0].Events, ID, termPosition)
	if err != nil {
		return errors.Wrap(err, "failed to update revindex word")
	}

	return nil
}

func (r eventRepository) updateRevindexWord(
	term string,
	eventIDsWithPositions []string,
	ID int64,
	termPosition int,
) error {
	eventIDWithPositionsList := append(eventIDsWithPositions, fmt.Sprintf("%d:%d", ID, termPosition))
	sort.Slice(eventIDWithPositionsList, func(i, j int) bool {
		first, _ := strconv.ParseInt(strings.Split(eventIDWithPositionsList[i], ":")[0], 10, 64)
		second, _ := strconv.ParseInt(strings.Split(eventIDWithPositionsList[j], ":")[0], 10, 64)

		return first < second
	})

	res := r.db.Model(&db_models.RevindexWord{}).
		Where("word = ?", term).
		Update("events", pq.StringArray(eventIDWithPositionsList))

	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to update revindex word in db")
	}

	return nil
}
