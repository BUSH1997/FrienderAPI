package revindex

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/goodsign/snowball"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

func (r eventRepository) Delete(ctx context.Context, event string, groupInfo models.GroupInfo) error {
	stemmer, err := snowball.NewWordStemmer("ru", "UTF_8")
	if err != nil {
		return errors.Wrap(err, "failed to init stammer")
	}

	defer stemmer.Close()

	existEvent, err := r.GetEventById(ctx, event)
	if err != nil {
		return errors.Wrapf(err, "failed to get event by uid %s", event)
	}

	dbRevindexEvent, err := r.getRevindexEvent(event)
	if err != nil {
		return errors.Wrapf(err, "failed to get revindex event by uid %s", event)
	}

	oldTerms, err := getTerms(strings.Split(existEvent.Title, " "), stemmer)
	if err != nil {
		return errors.Wrap(err, "failed to get old terms")
	}

	for _, term := range oldTerms {
		err := r.excludeEventIDFromRevindex(term, int64(dbRevindexEvent.ID))
		if err != nil {
			return errors.Wrapf(err, "failed to exclude event id from revindex of %s", term)
		}
	}

	err = r.events.Delete(ctx, event, groupInfo)
	if err != nil {
		return errors.Wrap(err, "failed to delete event")
	}

	return nil
}

func (r eventRepository) excludeEventIDFromRevindex(term string, ID int64) error {
	var dbRevindexWord db_models.RevindexWord

	res := r.db.Model(db_models.RevindexWord{}).
		Where("word = ?", term).
		Take(&dbRevindexWord)

	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to get old term")
	}

	if len(dbRevindexWord.Events) == 0 {
		return errors.New("expected at least one event id in term revindex")
	}

	eventIDPosition := getEventIDPosition(dbRevindexWord.Events, ID)
	if eventIDPosition == -1 {
		return errors.New("expected event id in term revindex")
	}

	var eventIDList pq.Int64Array
	if len(dbRevindexWord.Events) > 1 {
		eventIDList = append(
			dbRevindexWord.Events[0:eventIDPosition],
			dbRevindexWord.Events[eventIDPosition+1:len(dbRevindexWord.Events)]...,
		)
	}

	sort.Slice(eventIDList, func(i, j int) bool {
		return eventIDList[i] < eventIDList[j]
	})

	res = r.db.Model(&db_models.RevindexWord{}).
		Where("word = ?", term).
		Update("events", eventIDList)
	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to update old term")
	}

	return nil
}

func (r eventRepository) getRevindexEvent(uid string) (db_models.RevindexEvent, error) {
	var dbRevindexEvent db_models.RevindexEvent
	res := r.db.Model(db_models.RevindexEvent{}).
		Where("uid = ?", uid).
		Take(&dbRevindexEvent)
	if err := res.Error; err != nil {
		return db_models.RevindexEvent{}, errors.Wrap(err, "failed to get revindex event")
	}

	return dbRevindexEvent, nil
}
