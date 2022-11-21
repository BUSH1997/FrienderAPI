package user

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type UseCase interface {
	GetRefresh(ctx context.Context) (models.RefreshToken, error)
	CheckRefresh(ctx context.Context, old models.RefreshToken, new models.RefreshToken) error
	UpdateRefresh(ctx context.Context, fingerPrint models.FingerPrintData) (models.RefreshToken, error)
	GenerateAuthToken(ctx context.Context) (models.AuthToken, error)
}
