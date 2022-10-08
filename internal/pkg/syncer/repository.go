package syncer

import (
	"context"
	"time"
)

type Repository interface {
	GetUpdatedTime(ctx context.Context) (time.Time, error)
	Update(ctx context.Context, time time.Time) error
}
