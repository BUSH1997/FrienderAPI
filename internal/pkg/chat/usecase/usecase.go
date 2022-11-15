package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/chat"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
)

type ChatUsecase struct {
	chatRepository chat.Repository
	logger         hardlogger.Logger
}

func New(chatRepository chat.Repository, logger hardlogger.Logger) chat.Usecase {
	return &ChatUsecase{
		chatRepository: chatRepository,
		logger:         logger,
	}
}
