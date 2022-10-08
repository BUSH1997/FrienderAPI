package syncer

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/syncer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type SyncerConfig struct {
	SyncOffset           time.Duration `mapstructure:"sync_offset"`
	WaitTime             time.Duration `mapstructure:"wait_time"`
	ResyncDelayAfterFail time.Duration `mapstructure:"resync_delay_after_fail"`
	URL                  string        `mapstructure:"url"`
}

type PublicSyncer struct {
	Logger             *logrus.Logger
	Syncer             SyncerConfig
	PublicEventsClient client.PublicEventsClient
	Events             event.Usecase
	Repository         syncer.Repository
}

func New(
	config SyncerConfig,
	logger *logrus.Logger,
	client client.PublicEventsClient,
	usecase event.Usecase,
	repository syncer.Repository,
) PublicSyncer {
	return PublicSyncer{
		Logger:             logger,
		Syncer:             config,
		PublicEventsClient: client,
		Events:             usecase,
		Repository:         repository,
	}
}

func (s PublicSyncer) RunPublicSync() {
	ctx := context.Background()

	for {
		time.Sleep(s.Syncer.WaitTime)

		updatedAt, err := s.Repository.GetUpdatedTime(ctx)
		if err != nil {
			s.Logger.Warnf("failed to get updated_at with error: %s", err.Error())
			continue
		}

		if time.Since(updatedAt) < s.Syncer.SyncOffset {
			continue
		}

		err = s.syncPublicEvents(ctx)
		if err != nil {
			s.Logger.Warnf("failed to process sync with error: %s", err.Error())

			err = s.Repository.Update(ctx, time.Now().Add(-s.Syncer.ResyncDelayAfterFail))
			if err != nil {
				s.Logger.Warnf("failed to update sync after fail with error: %s", err.Error())
			}

			continue
		}

		err = s.Repository.Update(ctx, time.Now())
		if err != nil {
			s.Logger.Warnf("failed to update sync after success with error: %s", err.Error())
		}
	}
}

func (s PublicSyncer) syncPublicEvents(ctx context.Context) error {
	externalEvents, err := s.PublicEventsClient.UploadPublicEvents(ctx, s.Syncer.URL)
	if err != nil {
		return errors.Wrap(err, "failed to upload public events")
	}

	existingEvents, err := s.Events.GetAllPublic(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get existing events")
	}

	newEvents, changedEvents := getEventsToImport(existingEvents, externalEvents)

	for _, newEvent := range newEvents {
		_, err = s.Events.Create(ctx, newEvent)
		if err != nil {
			return errors.Wrapf(err, "failed to create public event, uid: %d", newEvent.Uid)
		}
	}

	for _, changedEvent := range changedEvents {
		err = s.Events.Update(ctx, changedEvent)
		if err != nil {
			return errors.Wrapf(err, "failed to update public event, uid: %d", changedEvent.Uid)
		}
	}

	s.Logger.Info("successfully synced public events")

	return nil
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
