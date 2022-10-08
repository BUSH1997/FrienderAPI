package postgres

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

func (r syncerRepository) Update(ctx context.Context, time time.Time) error {
	res := r.db.Model(&Syncer{}).Where("id = 1").Update("updated_at", time)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to update syncer time")
	}

	return nil
}
