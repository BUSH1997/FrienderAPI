package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type UserRepository struct {
	logger hardlogger.Logger
	db     *gorm.DB
}

func New(db *gorm.DB, logger hardlogger.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}
