package usecase

import (
	"context"
	"github.com/pkg/errors"
)

func (uc ChatUsecase) UpdateLastCheckTime(ctx context.Context, event string, user int64, time int64) error {
	err := uc.chatRepository.UpdateLastCheckTime(ctx, event, user, time)
	if err != nil {
		return errors.Wrap(err, "failed to update last check time")
	}

	return nil
}
