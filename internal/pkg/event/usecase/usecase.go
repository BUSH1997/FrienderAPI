package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/sirupsen/logrus"
)

type eventUsecase struct {
	Events event.Repository
	logger *logrus.Logger
}

func New(repository event.Repository, logger *logrus.Logger) event.Usecase {
	return &eventUsecase{
		Events: repository,
		logger: logger,
	}
}
