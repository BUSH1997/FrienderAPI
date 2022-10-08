package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (u eventUsecase) Create(ctx context.Context, event models.Event) error {
	err := u.Events.Create(ctx, event)
	if err != nil {
		return errors.Wrap(err, "failed to create public event in usecase")
	}

	return nil
}
