package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/chat"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type chatRepository struct {
	db     *gorm.DB
	logger hardlogger.Logger
}

func New(db *gorm.DB, logger hardlogger.Logger) chat.Repository {
	return &chatRepository{
		db:     db,
		logger: logger,
	}
}
