package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func (uc eventUsecase) Create(ctx context.Context, event models.Event) (models.Event, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return models.Event{}, errors.Wrap(err, "failed to  generate uuid")
	}

	if !event.IsPublic {
		event.Uid = uid.String()
	}

	if event.IsPublic {
		event.GeoData = models.Geo{
			Longitude: 127,
			Latitude:  127,
		}
	}

	//err = uc.validateEvent(event)
	//if err != nil {
	//	return models.Event{}, errors.Wrap(err, " failed to validate event")
	//}

	err = uc.Events.Create(ctx, event)
	if err != nil {
		return models.Event{}, errors.Wrap(err, "failed to create public event in usecase")
	}

	return event, nil
}
