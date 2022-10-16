package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type profileRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) profile.Repository {
	return &profileRepository{
		db:     db,
		logger: logger,
	}
}
