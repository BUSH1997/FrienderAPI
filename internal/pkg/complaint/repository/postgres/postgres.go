package postgres

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/complaint"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"gorm.io/gorm"
)

type complaintRepository struct {
	db     *gorm.DB
	logger hardlogger.Logger
}

func New(db *gorm.DB, logger hardlogger.Logger) complaint.Repository {
	return &complaintRepository{
		db:     db,
		logger: logger,
	}
}
