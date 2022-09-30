package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/syncer"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type syncerRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) syncer.Repository {
	return syncerRepository{
		db:     db,
		logger: logger,
	}
}
