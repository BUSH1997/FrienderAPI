package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type eventRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) event.Repository {
	return eventRepository{
		db:     db,
		logger: logger,
	}
}
