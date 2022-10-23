package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
)

func (uc ChatUsecase) GetChats(ctx context.Context) ([]models.Chat, error) {
	chats, err := uc.chatRepository.GetChats(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chats in usecase")
	}

	return chats, nil
}

func (uc ChatUsecase) GetMessages(ctx context.Context, opts models.GetMessageOpts) ([]models.Message, error) {
	messages, err := uc.chatRepository.GetMessages(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get messages in usecase")
	}

	return messages, nil
}
