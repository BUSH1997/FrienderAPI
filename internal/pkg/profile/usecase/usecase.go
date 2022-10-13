package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/sirupsen/logrus"
)

type UseCase struct {
	Repository profile.Repository
	Logger *logrus.Logger
}

func New(repository profile.Repository, logger *logrus.Logger) profile.UseCase {
	return &UseCase{
		Repository: repository,
		Logger: logger,
	}
}