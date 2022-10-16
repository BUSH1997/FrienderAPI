package profile

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type Repository interface {
	UpdateProfile(ctx context.Context, profile models.ChangeProfile) error
}
