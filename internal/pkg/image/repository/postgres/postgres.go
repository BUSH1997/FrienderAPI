package postgres

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ImageRepositoryBD struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) ImageRepositoryBD {
	return ImageRepositoryBD{
		db:     db,
		logger: logger,
	}
}
