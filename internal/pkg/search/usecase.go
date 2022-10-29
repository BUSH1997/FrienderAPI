package search

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type UseCase interface {
	Search(ctx context.Context, searchData models.Search) ([]models.Event, error)
}
