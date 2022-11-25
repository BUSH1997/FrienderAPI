package usecase

import (
	"context"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"time"
)

func (u UserUseCase) GetRefresh(ctx context.Context) (models.RefreshToken, error) {
	token, err := u.UserRepo.GetRefresh(ctx)
	if err != nil {
		return models.RefreshToken{}, errors.Wrap(err, "failed to get refresh token in usecase")
	}

	return token, nil
}

func (u UserUseCase) CheckRefresh(ctx context.Context, old models.RefreshToken, new models.RefreshToken) error {
	if old.Value != new.Value {
		err := errors.New("invalid refresh token")
		u.Logger.WithError(err).Errorf("invalid refresh token")
		return err
	}

	if old.FingerPrint != new.FingerPrint {
		err := errors.New("invalid fingerprint")
		u.Logger.WithError(err).Errorf("invalid fingerprint")
		return err
	}

	return nil
}

func (u UserUseCase) GenerateAuthToken(ctx context.Context) (models.AuthToken, error) {
	user := contextlib.GetUser(ctx)

	expireTime := time.Now().Add(u.UserConfig.Cookie.Auth.Exp).Unix()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user,
		"exp": expireTime,
	})

	token, err := claims.SignedString([]byte(u.UserConfig.AuthSignSecret))
	if err != nil {
		return models.AuthToken{}, errors.Wrap(err, "failed to get signed token")
	}

	return models.AuthToken{
		UserID:  user,
		Value:   token,
		Expires: expireTime,
	}, nil
}
