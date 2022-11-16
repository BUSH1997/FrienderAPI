package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
)

func (uc ChatUsecase) UpdateLastCheckTime(ctx context.Context, event string, user int64, time int64) error {
	ctx = uc.logger.WithCaller(ctx)

	err := uc.chatRepository.UpdateLastCheckTime(ctx, event, user, time)
	if err != nil {
		return errors.Wrap(err, "failed to update last check time")
	}

	return nil
}
