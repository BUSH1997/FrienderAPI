package usecase

import "github.com/BUSH1997/FrienderAPI/internal/pkg/chat"

type ChatUsecase struct {
	chatRepository chat.Repository
}

func New(chatRepository chat.Repository) chat.Usecase {
	return &ChatUsecase{
		chatRepository: chatRepository,
	}
}
