package postgres

import (
	"context"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"time"
)

func (r syncerRepository) Update(ctx context.Context, time time.Time) error {
	res := r.db.Model(&db_models.Syncer{}).Where("id = 1").Update("updated_at", time)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to update syncer time")
	}

	return nil
}
