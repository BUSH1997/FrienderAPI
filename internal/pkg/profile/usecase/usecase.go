package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/award"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/status"
	"github.com/sirupsen/logrus"
)

type UseCase struct {
	profileRepository profile.Repository
	eventRepository   event.Repository
	awardRepository   award.Repository
	statusRepository  status.Repository
	Logger            *logrus.Logger
}

func New(
	profileRepository profile.Repository,
	eventRepository event.Repository,
	awardRepository award.Repository,
	statusRepository status.Repository,
	logger *logrus.Logger,
) profile.UseCase {
	return &UseCase{
		profileRepository: profileRepository,
		eventRepository:   eventRepository,
		awardRepository:   awardRepository,
		statusRepository:  statusRepository,
		Logger:            logger,
	}
}
