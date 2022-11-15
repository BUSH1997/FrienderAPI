package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/group"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type groupRepository struct {
	db     *gorm.DB
	logger hardlogger.Logger
}

func New(db *gorm.DB, logger hardlogger.Logger) group.Repository {
	return &groupRepository{
		db:     db,
		logger: logger,
	}
}
