package revindex

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/goodsign/snowball"
	"github.com/pkg/errors"
	"strings"
)

func (r eventRepository) Update(ctx context.Context, event models.Event) error {
	existEvent, err := r.GetEventById(ctx, event.Uid)
	if err != nil {
		return errors.Wrapf(err, "failed to get event by uid %s", event.Uid)
	}

	if existEvent.Title != event.Title {
		err := r.updateRevindex(event.Title, existEvent.Title, event.Uid)
		if err != nil {
			return errors.Wrapf(err, "failed to update revindex with event %s", event.Uid)
		}
	}

	err = r.events.Update(ctx, event)
	if err != nil {
		return errors.Wrap(err, "failed to update event")
	}

	return nil
}

func (r eventRepository) updateRevindex(newTitle string, oldTitle string, uid string) error {
	stemmer, err := snowball.NewWordStemmer("ru", "UTF_8")
	if err != nil {
		return errors.Wrap(err, "failed to init stammer")
	}

	defer stemmer.Close()

	dbRevindexEvent, err := r.getRevindexEvent(uid)
	if err != nil {
		return errors.Wrapf(err, "failed to get revindex event by uid %s", uid)
	}

	oldTerms, err := getTerms(strings.Split(oldTitle, " "), stemmer)
	if err != nil {
		return errors.Wrap(err, "failed to get old terms")
	}

	for _, term := range oldTerms {
		err := r.excludeEventIDFromRevindex(term, int64(dbRevindexEvent.ID))
		if err != nil {
			return errors.Wrapf(err, "failed to exclude event id from revindex of %s", term)
		}
	}

	newTerms, err := getTerms(strings.Split(newTitle, " "), stemmer)
	if err != nil {
		return errors.Wrap(err, "failed to get new terms")
	}

	for _, term := range newTerms {
		err := r.createOrUpdateRevindex(term, int64(dbRevindexEvent.ID))
		if err != nil {
			return errors.Wrapf(err, "failed to update revindex with event %s for term %s", uid, term)
		}
	}

	return nil
}

func getEventIDPosition(eventIDs []int64, ID int64) int {
	for i, oldEventID := range eventIDs {
		if oldEventID == ID {
			return i
		}
	}

	return -1
}
