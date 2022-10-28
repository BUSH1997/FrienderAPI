package revindex

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type eventRepository struct {
	db     *gorm.DB
	events event.Repository
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger, eventRepo event.Repository) event.Repository {
	return &eventRepository{
		db:     db,
		events: eventRepo,
		logger: logger,
	}
}
