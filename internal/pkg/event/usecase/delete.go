package usecase

import (
	"context"
	"github.com/pkg/errors"
)

func (uc eventUsecase) DeletePhotos(ctx context.Context, links []string, uid string) error {
	err := uc.Events.DeletePhotos(ctx, links, uid)
	if err != nil {
		return errors.Wrap(err, "failed to delete event photos in usecase")
	}

	return nil
}
