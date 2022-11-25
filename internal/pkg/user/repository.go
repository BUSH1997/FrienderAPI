package user

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	GetRefresh(ctx context.Context) (models.RefreshToken, error)
	UpdateRefresh(ctx context.Context, token models.RefreshToken) error
}
