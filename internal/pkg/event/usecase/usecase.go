package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/blacklist"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/sirupsen/logrus"
)

type eventUsecase struct {
	Events      event.Repository
	BlackLister blacklist.BlackLister
	logger      *logrus.Logger
}

func New(repository event.Repository, blackLister blacklist.BlackLister, logger *logrus.Logger) event.Usecase {
	return &eventUsecase{
		Events:      repository,
		BlackLister: blackLister,
		logger:      logger,
	}
}
