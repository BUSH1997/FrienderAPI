package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/search"
	"github.com/sirupsen/logrus"
)

type UseCase struct {
	searchRepository search.Repository
	eventRepository  event.Repository
	skipList         map[string]bool
	logger           *logrus.Logger
}

func New(
	searchRepository search.Repository,
	eventRepository event.Repository,
	skipList []string,
	logger *logrus.Logger,
) search.UseCase {
	skipMap := make(map[string]bool)
	for _, skip := range skipList {
		skipMap[skip] = true
	}

	return &UseCase{
		searchRepository: searchRepository,
		eventRepository:  eventRepository,
		skipList:         skipMap,
		logger:           logger,
	}
}
