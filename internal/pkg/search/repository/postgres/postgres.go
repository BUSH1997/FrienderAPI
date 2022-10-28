package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/search"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type searchRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) search.Repository {
	return &searchRepository{
		db:     db,
		logger: logger,
	}
}
