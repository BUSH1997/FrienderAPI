package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
)

func (uc ChatUsecase) DeleteMessage(ctx context.Context, messageID string) error {
	ctx = uc.logger.WithCaller(ctx)

	err := uc.chatRepository.DeleteMessage(ctx, messageID)
	if err != nil {
		return errors.Wrap(err, "failed to delete message in usecase")
	}

	return nil
}
