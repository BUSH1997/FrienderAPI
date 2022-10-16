package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/status"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type statusRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) status.Repository {
	return &statusRepository{
		db:     db,
		logger: logger,
	}
}
