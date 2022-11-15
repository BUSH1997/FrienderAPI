package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/blacklist"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/search"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
)

type eventUsecase struct {
	Events            event.Repository
	ProfileRepository profile.Repository
	SearchRepository  search.Repository
	BlackLister       blacklist.BlackLister
	skipList          map[string]bool
	logger            hardlogger.Logger
}

func New(
	repository event.Repository,
	profileRepository profile.Repository,
	searchRepository search.Repository,
	blackLister blacklist.BlackLister,
	skipList []string,
	logger hardlogger.Logger,
) event.Usecase {
	skipMap := make(map[string]bool)
	for _, skip := range skipList {
		skipMap[skip] = true
	}

	return &eventUsecase{
		Events:            repository,
		ProfileRepository: profileRepository,
		SearchRepository:  searchRepository,
		BlackLister:       blackLister,
		skipList:          skipMap,
		logger:            logger,
	}
}
