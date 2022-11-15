package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type profileRepository struct {
	db     *gorm.DB
	logger hardlogger.Logger
}

func New(db *gorm.DB, logger hardlogger.Logger) profile.Repository {
	return &profileRepository{
		db:     db,
		logger: logger,
	}
}
