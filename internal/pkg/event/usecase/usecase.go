package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/blacklist"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/sirupsen/logrus"
)

type eventUsecase struct {
	Events            event.Repository
	BlackLister       blacklist.BlackLister
	logger            *logrus.Logger
	ProfileRepository profile.Repository
}

func New(repository event.Repository, ProfileRepository profile.Repository, blackLister blacklist.BlackLister, logger *logrus.Logger) event.Usecase {
	return &eventUsecase{
		Events:            repository,
		BlackLister:       blackLister,
		logger:            logger,
		ProfileRepository: ProfileRepository,
	}
}
