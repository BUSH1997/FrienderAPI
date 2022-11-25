package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/algorithms"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/google/uuid"
	"time"
)

func (u UserUseCase) UpdateRefresh(ctx context.Context, fingerPrint models.FingerPrintData) (models.RefreshToken, error) {
	refreshToken := models.RefreshToken{
		Value:       uuid.New().String(),
		Expires:     time.Now().Add(u.UserConfig.Cookie.Refresh.Exp).Unix(),
		FingerPrint: algorithms.GetFingerPrint([]string{fingerPrint.UserAgent, fingerPrint.UserIP}),
	}

	err := u.UserRepo.UpdateRefresh(ctx, refreshToken)
	if err != nil {
		return models.RefreshToken{}, errors.Wrap(err, "failed to update refresh token in usecase")
	}

	return refreshToken, nil
}
