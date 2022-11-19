package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	httplib "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
)

type UseCase struct {
	profileRepository profile.Repository
	eventRepository   event.Repository
	Logger            hardlogger.Logger
	httpClient        httplib.Client
}

func New(
	profileRepository profile.Repository,
	eventRepository event.Repository,
	logger hardlogger.Logger,
	client httplib.Client,
) profile.UseCase {
	return &UseCase{
		profileRepository: profileRepository,
		eventRepository:   eventRepository,
		Logger:            logger,
		httpClient:        client,
	}
}
