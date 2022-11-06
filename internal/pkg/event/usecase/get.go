package usecase

import (
	"context"
	context2 "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/stammer"
	"github.com/pkg/errors"
	"sort"
	"strings"
	"time"
)

func (uc eventUsecase) GetAllPublic(ctx context.Context) ([]models.Event, error) {
	events, err := uc.Events.GetAll(ctx, models.GetEventParams{
		IsPublic: models.DefinedBool(true),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all public events in usecase")
	}

	return events, nil
}

func (uc eventUsecase) GetEventById(ctx context.Context, id string) (models.Event, error) {
	return uc.Events.GetEventById(ctx, id)
}

func (uc eventUsecase) GetAllCategories(ctx context.Context) ([]string, error) {
	return uc.Events.GetAllCategories(ctx)
}

func (uc eventUsecase) Get(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	events, err := uc.routerGet(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get events in usecase")
	}

	sort.SliceStable(events, func(i, j int) bool {
		return events[i].StartsAt < events[j].StartsAt
	})

	if params.SortMembers != "" {
		sort.SliceStable(events, func(i, j int) bool {
			if params.SortMembers == "asc" {
				return len(events[i].Members) > len(events[j].Members)
			}

			return len(events[i].Members) < len(events[j].Members)
		})
	}

	if events == nil {
		events = make([]models.Event, 0)
	}

	return events, nil
}

func (uc eventUsecase) GetSubscribeEvent(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	subscribes, err := uc.ProfileRepository.GetSubscribe(ctx, context2.GetUser(ctx))
	if err != nil {
		uc.logger.WithError(err).Errorf("[GetSubscribeEvent]")
		return []models.Event{}, err
	}

	var result []models.Event
	for _, subscribe := range subscribes {
		var profileEvents []models.Event
		if subscribe.IsGroup {
			params.IsActive = models.Bool{Defined: true, Value: true}
			params.GroupId = subscribe.Id

			profileEvents, err = uc.Events.GetGroupEvent(ctx, params)
			if err != nil {
				uc.logger.WithError(err).Errorf("[GetSubscribeEvent] faile getgroup event")
				return []models.Event{}, err
			}
		} else {
			profileEvents, err = uc.Events.GetSharings(ctx, models.GetEventParams{
				UserID:   subscribe.Id,
				IsActive: models.DefinedBool(true),
			})
			if err != nil {
				uc.logger.WithError(err).Errorf("[GetSubscribeEvent] faile getgroup event")
				return []models.Event{}, err
			}
		}

		result = append(result, profileEvents...)
	}

	return result, nil
}

func (uc eventUsecase) routerGet(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	if params.IsActive.IsDefined() && params.UserID != 0 {
		return uc.Events.GetSharings(ctx, params)
	}

	if params.IsSubscriber.IsDefinedTrue() {
		return uc.Events.GetSubscriptionEvents(ctx, params.UserID)
	}
	if params.GroupId != 0 && (params.IsAdmin.IsDefinedTrue() || params.IsAdmin.IsDefinedFalse()) {
		return uc.Events.GetGroupAdminEvent(ctx, params)
	}
	if params.GroupId != 0 {
		return uc.Events.GetGroupEvent(ctx, params)
	}
	if params.Source == "subscribe" {
		return uc.GetSubscribeEvent(ctx, params)
	}
	if params.Search.Enabled {
		return uc.GetSearch(ctx, params)
	}

	return uc.Events.GetAll(ctx, params)
}

func (uc eventUsecase) GetSearch(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	stammers, err := stammer.GetStammers(stammer.FilterSkipList(params.Search.SearchData.Words, uc.skipList))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stammers")
	}

	eventUIDs, err := uc.SearchRepository.GetEventUIDs(ctx, stammers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get event uids")
	}

	params.UIDs = eventUIDs

	events, err := uc.Events.GetAll(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get events by uids")
	}

	eventRatesMap := make(map[string]float64)
	for _, event := range events {
		eventRatesMap[event.Uid] = float64(len(stammers)) / float64(len(strings.Split(event.Title, " ")))
	}

	sort.Slice(events, func(i, j int) bool {
		return eventRatesMap[events[i].Uid] > eventRatesMap[events[j].Uid]
	})

	events = Filter(events, params.Search.SearchData.Sources)

	return events, nil
}

func Filter(events []models.Event, sources []string) []models.Event {
	if len(sources) == 0 {
		return events
	}

	sourceMap := make(map[string]bool)
	for _, source := range sources {
		sourceMap[source] = true
	}

	ret := make([]models.Event, 0, len(events))
	for _, event := range events {
		if sourceMap[event.Source] && event.StartsAt > time.Now().Unix() {
			ret = append(ret, event)
		}
	}

	return ret
}
