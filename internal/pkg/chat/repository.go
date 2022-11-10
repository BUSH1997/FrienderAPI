package chat

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	CreateMessage(ctx context.Context, message models.Message) error
	GetMessages(ctx context.Context, opts models.GetMessageOpts) ([]models.Message, error)
	GetChats(ctx context.Context) ([]models.Chat, error)
	UpdateLastCheckTime(ctx context.Context, event string, user int64, time int64) error
}
