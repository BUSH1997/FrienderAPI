package http

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/user"
)

type UserHandler struct {
	Logger      hardlogger.Logger
	UserUseCase user.UseCase
	AuthSecret  string
}

func NewUserHandler(usecase user.UseCase, authSecret string, logger hardlogger.Logger) *UserHandler {
	return &UserHandler{
		AuthSecret:  authSecret,
		UserUseCase: usecase,
		Logger:      logger,
	}
}
