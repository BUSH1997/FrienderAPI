package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/search"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type searchRepository struct {
	db     *gorm.DB
	logger hardlogger.Logger
}

func New(db *gorm.DB, logger hardlogger.Logger) search.Repository {
	return &searchRepository{
		db:     db,
		logger: logger,
	}
}
