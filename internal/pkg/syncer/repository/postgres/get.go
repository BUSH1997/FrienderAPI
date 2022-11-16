package postgres

import (
	"context"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"time"
)

func (r syncerRepository) GetUpdatedTime(ctx context.Context) (time.Time, error) {
	var syncer db_models.Syncer

	res := r.db.Take(&syncer)
	if err := res.Error; err != nil {
		return time.Time{}, errors.Wrap(err, "failed to get updated time")
	}

	return syncer.UpdatedAt, nil
}
