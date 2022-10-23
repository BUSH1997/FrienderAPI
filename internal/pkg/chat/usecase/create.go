package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

func (uc ChatUsecase) CreateMessage(ctx context.Context, message models.Message) error {
	return uc.chatRepository.CreateMessage(ctx, message)
}
