package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/status"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type statusRepository struct {
	db     *gorm.DB
	logger hardlogger.Logger
}

func New(db *gorm.DB, logger hardlogger.Logger) status.Repository {
	return &statusRepository{
		db:     db,
		logger: logger,
	}
}
