package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/syncer"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type syncerRepository struct {
	db     *gorm.DB
	logger hardlogger.Logger
}

func New(db *gorm.DB, logger hardlogger.Logger) syncer.Repository {
	return syncerRepository{
		db:     db,
		logger: logger,
	}
}
