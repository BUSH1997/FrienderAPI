package postgres

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

func (r syncerRepository) GetUpdatedTime(ctx context.Context) (time.Time, error) {
	var syncer Syncer

	res := r.db.Take(&syncer)
	if err := res.Error; err != nil {
		return time.Time{}, errors.Wrap(err, "failed to get updated time")
	}

	return syncer.UpdatedAt, nil
}
