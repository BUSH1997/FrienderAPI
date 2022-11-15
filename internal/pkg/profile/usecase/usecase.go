package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/award"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/status"
	httplib "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
)

type UseCase struct {
	profileRepository profile.Repository
	eventRepository   event.Repository
	awardRepository   award.Repository
	statusRepository  status.Repository
	Logger            hardlogger.Logger
	httpClient        httplib.Client
}

func New(
	profileRepository profile.Repository,
	eventRepository event.Repository,
	awardRepository award.Repository,
	statusRepository status.Repository,
	logger hardlogger.Logger,
	client httplib.Client,
) profile.UseCase {
	return &UseCase{
		profileRepository: profileRepository,
		eventRepository:   eventRepository,
		awardRepository:   awardRepository,
		statusRepository:  statusRepository,
		Logger:            logger,
		httpClient:        client,
	}
}
