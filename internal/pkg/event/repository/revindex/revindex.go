package revindex

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type eventRepository struct {
	db       *gorm.DB
	events   event.Repository
	skipList map[string]bool
	logger   hardlogger.Logger
}

func New(db *gorm.DB, logger hardlogger.Logger, eventRepo event.Repository, skipList []string) event.Repository {
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
