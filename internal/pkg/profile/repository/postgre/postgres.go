package postgre

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) profile.Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}
