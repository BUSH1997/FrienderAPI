package usecase

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/search"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UseCase struct {
	searchRepository search.Repository
	eventRepository  event.Repository
	logger           *logrus.Logger
}

func New(
	searchRepository search.Repository,
	eventRepository event.Repository,
	logger *logrus.Logger,
) search.UseCase {
	return &UseCase{
		searchRepository: searchRepository,
		eventRepository:  eventRepository,
		logger:           logger,
	}
}

func (uc UseCase) Search(words []string) ([]models.Event, error) {
	events, err := uc.searchRepository.Search(words)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get events in search")
	}

	return events, nil
}
