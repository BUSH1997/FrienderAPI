package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/group"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type groupRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func New(db *gorm.DB, logger *logrus.Logger) group.Repository {
	return &groupRepository{
		db:     db,
		logger: logger,
	}
}
