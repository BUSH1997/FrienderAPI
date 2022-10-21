package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

func (gu *groupUseCase) Create(ctx context.Context, group models.Group) error {
	return gu.repository.Create(ctx, group)
}
