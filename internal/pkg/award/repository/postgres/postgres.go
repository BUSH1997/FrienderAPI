package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/award"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type awardRepository struct {
	db     *gorm.DB
	logger hardlogger.Logger
}

func New(db *gorm.DB, logger hardlogger.Logger) award.Repository {
	return &awardRepository{
		db:     db,
		logger: logger,
	}
}
