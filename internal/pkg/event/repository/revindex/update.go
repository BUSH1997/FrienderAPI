package revindex

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/stammer"
	"strings"
)

func (r eventRepository) Update(ctx context.Context, event models.Event) error {
	ctx = r.logger.WithCaller(ctx)

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
	dbRevindexEvent, err := r.getRevindexEvent(uid)
	if err != nil {
		return errors.Wrapf(err, "failed to get revindex event by uid %s", uid)
	}

	oldTerms, err := stammer.GetStammers(stammer.FilterSkipList(strings.Split(oldTitle, " "), r.skipList))
	if err != nil {
		return errors.Wrap(err, "failed to get stammers from title")
	}

	for _, term := range oldTerms {
		err := r.excludeEventIDFromRevindex(term, int64(dbRevindexEvent.ID))
		if err != nil {
			return errors.Wrapf(err, "failed to exclude event id from revindex of %s", term)
		}
	}

	newTerms, err := stammer.GetStammers(stammer.FilterSkipList(strings.Split(newTitle, " "), r.skipList))
	if err != nil {
		return errors.Wrap(err, "failed to get stammers from title")
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
