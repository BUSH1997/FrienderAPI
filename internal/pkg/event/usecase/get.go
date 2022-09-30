package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (u eventUsecase) GetAllPublic(ctx context.Context) ([]models.Event, error) {
	events, err := u.Events.GetAllPublic(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all public events in usecase")
	}

	return events, nil
}
