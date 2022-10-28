package revindex

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type eventRepository struct {
	db       *gorm.DB
	events   event.Repository
	skipList map[string]bool
	logger   *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger, eventRepo event.Repository, skipList []string) event.Repository {
	skipMap := make(map[string]bool)
	for _, skip := range skipList {
		skipMap[skip] = true
	}

	return &eventRepository{
		db:       db,
		events:   eventRepo,
		skipList: skipMap,
		logger:   logger,
	}
}
