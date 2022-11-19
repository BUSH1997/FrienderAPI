package postgres

import (
	"context"
	complaint_pkg "github.com/BUSH1997/FrienderAPI/internal/pkg/complaint"
	context2 "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"strconv"
	"time"
)

func (r complaintRepository) Create(ctx context.Context, complaint models.Complaint) error {
	dbComplaint := db_models.Complaint{
		TimeCreated: time.Now().Unix(),
		IsProcessed: false,
	}

	var dbInitiator db_models.User
	res := r.db.Take(&dbInitiator, "uid = ?", context2.GetUser(ctx))
	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to get initiator")
	}

	dbComplaint.Initiator = int64(dbInitiator.ID)

	if complaint.Event != "" {
		dbComplaint.Item = "event"
		dbComplaint.ItemUID = complaint.Event
	}

	if complaint.User != 0 {
		dbComplaint.Item = "user"
		dbComplaint.ItemUID = strconv.Itoa(int(complaint.User))
	}

	res = r.db.Create(&dbComplaint)
	if err := res.Error; err != nil {
		if postgres.ProcessError(err) == postgres.UniqueViolationError {
			return errors.Transform(err, complaint_pkg.ErrAlreadyExists)
		}

		return errors.Wrap(err, "failed to create complaint")
	}

	return nil
}
