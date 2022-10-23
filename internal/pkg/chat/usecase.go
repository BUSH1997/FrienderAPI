package chat

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Usecase interface {
	CreateMessage(ctx context.Context, message models.Message) error
	GetChats(ctx context.Context) ([]models.Chat, error)
	GetMessages(ctx context.Context, opts models.GetMessageOpts) ([]models.Message, error)
}
