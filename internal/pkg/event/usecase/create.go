package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/gofrs/uuid"
)

func (uc eventUsecase) Create(ctx context.Context, event models.Event) (models.Event, error) {
	ctx = uc.logger.WithCaller(ctx)

	uid, err := uuid.NewV4()
	if err != nil {
		return models.Event{}, errors.Wrap(err, "failed to  generate uuid")
	}

	if event.Source != "vk_event" {
		event.Uid = uid.String()
	}

	if event.GroupInfo.GroupId != 0 {
		event.IsPublic = true
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
