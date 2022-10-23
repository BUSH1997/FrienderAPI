package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/chat"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type chatRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) chat.Repository {
	return &chatRepository{
		db:     db,
		logger: logger,
	}
}
