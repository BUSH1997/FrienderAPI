package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/award"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type awardRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) award.Repository {
	return &awardRepository{
		db:     db,
		logger: logger,
	}
}
