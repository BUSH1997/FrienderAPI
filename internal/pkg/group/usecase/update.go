package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (gu *groupUseCase) Update(ctx context.Context, group models.Group) error {
	err := gu.repository.Update(ctx, group)
	if err != nil {
		return errors.Wrap(err, "failed to update group in usecase")
	}

	return nil
}
