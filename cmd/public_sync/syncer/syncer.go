package syncer

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client"
	context2 "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/syncer"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"strconv"
	"strings"
	"time"
)

type SyncerConfig struct {
	SyncOffset           time.Duration `mapstructure:"sync_offset"`
	WaitTime             time.Duration `mapstructure:"wait_time"`
	ResyncDelayAfterFail time.Duration `mapstructure:"resync_delay_after_fail"`
	Timepad              Timepad       `mapstructure:"timepad"`
	VK                   VK            `mapstructure:"vk"`
}

type Timepad struct {
	URL string `mapstructure:"url"`
}

type VK struct {
	GetEventsURL     string `mapstructure:"get_events"`
	GetEventsDataURL string `mapstructure:"get_events_data"`
}

type Syncer interface {
	SyncData() client.SyncData
	Client() client.PublicEventsClient
}

type SyncManager struct {
	Config       SyncerConfig
	Logger       hardlogger.Logger
	Syncers      []Syncer
	SyncerUpdate Syncer
	Events       event.Usecase
	Repository   syncer.Repository
}

func New(
	config SyncerConfig,
	logger hardlogger.Logger,
	syncers []Syncer,
	syncerUpdate Syncer,
	usecase event.Usecase,
	repository syncer.Repository,
) SyncManager {
	return SyncManager{
		Config:       config,
		Logger:       logger,
		Syncers:      syncers,
		SyncerUpdate: syncerUpdate,
		Events:       usecase,
		Repository:   repository,
	}
}

func (s SyncManager) RunPublicSync() {
	ctx := context.Background()

	for {
		time.Sleep(s.Config.WaitTime)

		updatedAt, err := s.Repository.GetUpdatedTime(ctx)
		if err != nil {
			s.Logger.Errorf("failed to get updated_at with error: %s", err.Error())
			continue
		}

		if time.Since(updatedAt) < s.Config.SyncOffset {
			continue
		}

		err = s.syncPublicEvents(ctx)
		if err != nil {
			s.Logger.Errorf("failed to process sync with error: %s", err.Error())

			err = s.Repository.Update(ctx, time.Now().Add(-s.Config.ResyncDelayAfterFail))
			if err != nil {
				s.Logger.Errorf("failed to update sync after fail with error: %s", err.Error())
			}
		}

		err = s.updatePublicEvents(ctx)
		if err != nil {
			s.Logger.Errorf("failed to process updateEvents sync with error: %s", err.Error())

			err = s.Repository.Update(ctx, time.Now().Add(-s.Config.ResyncDelayAfterFail))
			if err != nil {
				s.Logger.Errorf("failed to update sync after fail with error: %s", err.Error())
			}

			continue
		}

		err = s.Repository.Update(ctx, time.Now())
		if err != nil {
			s.Logger.Errorf("failed to update sync after success with error: %s", err.Error())
		}
	}
}

func (s SyncManager) syncPublicEvents(ctx context.Context) error {
	for _, syncer := range s.Syncers {
		ctx = context2.SetUser(ctx, 1) // TODO: authorize users for public syncers(in config)

		countItem, err := syncer.Client().GetCountPublicEventsWithSyncData(ctx, syncer.SyncData())
		if err != nil {
			return errors.Wrap(err, "failed to get count event")
		}
		var currentItem int = 0

		for currentItem < countItem && currentItem < 1000 {
			syncData := syncer.SyncData()
			syncData.GetFormData()[0]["offset"] = strconv.Itoa(currentItem)

			externalEvents, err := syncer.Client().UploadPublicEvents(ctx, syncer.SyncData())
			if err != nil {
				return errors.Wrap(err, "failed to upload public events")
			}

			newEvents, err := s.filterExistEvents(ctx, externalEvents)

			for _, newEvent := range newEvents {
				if err != nil {
					return errors.Wrapf(err, "failed to build blackList")
				}
				if strings.Contains(strings.ToLower(newEvent.Title), "отме") {
					continue
				}

				if newEvent.StartsAt > time.Now().Unix() && newEvent.GeoData.Address != "" {
					_, err = s.Events.Create(ctx, newEvent)
					if err != nil {
						return errors.Wrapf(err, "failed to create public event, uid: %d", newEvent.Uid)
					}
				}
			}

			currentItem += 100
		}

		s.Logger.Info("successfully synced public events")
	}

	return nil
}

func (s SyncManager) updatePublicEvents(ctx context.Context) error {
	countEvents, err := s.Events.GetCountEvents(ctx, models.SOURCE_EVENT_VK)
	if err != nil {
		return err
	}

	for i := int64(0); i < countEvents; i = i + 100 {
		param := models.GetEventParams{
			Source: models.SOURCE_EVENT_VK,
			Page:   int(i / 100),
			Limit:  100,
		}

		events, err := s.Events.Get(ctx, param)
		if err != nil {
			return err
		}

		var uuids []string
		for _, value := range events {
			uuids = append(uuids, value.Uid)
		}

		s.SyncerUpdate.Client().GetGroupsByIds(ctx, s.SyncerUpdate.SyncData(), uuids)
	}

	return nil
}

func (s SyncManager) filterExistEvents(ctx context.Context, externalEvents []models.Event) ([]models.Event, error) {
	var newEvents []models.Event
	for _, v := range externalEvents {
		isExist, err := s.Events.CheckIfExists(ctx, v)
		if err != nil {
			return []models.Event{}, err
		}
		if !isExist {
			newEvents = append(newEvents, v)
		}
	}

	return newEvents, nil
}

func getEventsToImport(existingEvents []models.Event, externalEvents []models.Event) ([]models.Event, []models.Event) {
	newEvents, oldEvents := separateNewAndOldEvents(existingEvents, externalEvents)

	changedEvents := getChangedEvents(existingEvents, oldEvents)

	return newEvents, changedEvents
}

func separateNewAndOldEvents(
	existingEvents []models.Event,
	externalEvents []models.Event,
) ([]models.Event, []models.Event) {
	existingEventsMap := make(map[string]bool)
	for _, existingEvent := range existingEvents {
		existingEventsMap[existingEvent.Uid] = true
	}

	newEvents := make([]models.Event, 0, len(externalEvents))
	oldEvents := make([]models.Event, 0, len(existingEvents))

	for _, externalEvent := range externalEvents {
		if existingEventsMap[externalEvent.Uid] {
			oldEvents = append(oldEvents, externalEvent)
			continue
		}

		newEvents = append(newEvents, externalEvent)
	}

	return newEvents, oldEvents
}

func getChangedEvents(existingEvents []models.Event, externalEvents []models.Event) []models.Event {
	existingEventsMap := make(map[string]string)
	for _, existingEvent := range existingEvents {
		existingEventsMap[existingEvent.Uid] = existingEvent.GetEtag()
	}

	changedEvents := make([]models.Event, 0, len(externalEvents))

	for _, externalEvent := range externalEvents {
		if existingEventsMap[externalEvent.Uid] == externalEvent.GetEtag() {
			continue
		}

		changedEvents = append(changedEvents, externalEvent)
	}

	return changedEvents
}
