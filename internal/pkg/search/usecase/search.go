package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/stammer"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

func (uc UseCase) Search(ctx context.Context, searchData models.Search) ([]models.Event, error) {
	stammers, err := stammer.GetStammers(stammer.FilterSkipList(searchData.Words, uc.skipList))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stammers")
	}

	eventUIDs, err := uc.searchRepository.GetEventUIDs(ctx, stammers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get event uids")
	}

	eventRatesMap := make(map[string]float64)

	events := make([]models.Event, 0, len(eventUIDs))
	for _, eventUID := range eventUIDs {
		event, err := uc.eventRepository.GetEventById(ctx, eventUID)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get event by uid %s", eventUID)
		}

		events = append(events, event)
	}

	for _, event := range events {
		eventRatesMap[event.Uid] = float64(len(stammers)) / float64(len(strings.Split(event.Title, " ")))
	}

	sort.Slice(events, func(i, j int) bool {
		return eventRatesMap[events[i].Uid] > eventRatesMap[events[j].Uid]
	})

	events = FilterBySource(events, searchData.Source)

	return events, nil
}

func FilterBySource(events []models.Event, source string) []models.Event {
	ret := make([]models.Event, 0, len(events))
	for _, event := range events {
		if event.Source != source {
			continue
		}

		ret = append(ret, event)
	}

	return ret
}
