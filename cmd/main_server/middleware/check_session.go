package middleware

import (
	"fmt"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/algorithms"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/user"
	userUsecase "github.com/BUSH1997/FrienderAPI/internal/pkg/user/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

func CheckSession(logger hardlogger.Logger, useCase user.UseCase, authConfig userUsecase.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			if strings.Contains(echoCtx.Path(), "auth") {
				return next(echoCtx)
			}

			ctx := echoCtx.Request().Context()
			user := contextlib.GetUser(ctx)

			authToken, err := getAuthToken(echoCtx, authConfig, logger)
			if err != nil {
				logger.WithError(err).Errorf("failed to get auth token")
				return echoCtx.JSON(http.StatusUnauthorized, err.Error())
			}

			if authToken != (models.AuthToken{}) {
				if authToken.UserID != user {
					err := errors.New("users mismatch")
					logger.WithError(err).Errorf("users mismatch %d and %d", authToken.UserID, user)
					return echoCtx.JSON(http.StatusUnauthorized, err.Error())
				}

				return next(echoCtx)
			}

			refreshCookie, err := getCookie(echoCtx, "refresh")
			if err != nil {
				return errors.Wrap(err, "failed to get refresh cookie")
			}

			if refreshCookie == nil {
				err := errors.New("empty refresh token")
				logger.WithError(err).Errorf("failed to auth user")
				return echoCtx.JSON(http.StatusBadRequest, err.Error())
			}

			userAgent := echoCtx.Request().Header.Get("User-Agent")
			userIP := echoCtx.Request().Header.Get("X-Forwarded-For")

			refreshToken := models.RefreshToken{
				Value:       refreshCookie.Value,
				Expires:     refreshCookie.Expires.Unix(),
				FingerPrint: algorithms.GetFingerPrint([]string{userAgent, userIP}),
			}

			refreshFromDB, err := useCase.GetRefresh(ctx)
			if err != nil {
				return echoCtx.JSON(http.StatusUnauthorized, errors.Wrap(err, "failed to get refresh token").Error())
			}

			err = useCase.CheckRefresh(ctx, refreshFromDB, refreshToken)
			if err != nil {
				return echoCtx.JSON(http.StatusUnauthorized, errors.Wrap(err, "failed to check refresh token").Error())
			}

			newAuthToken, err := useCase.GenerateAuthToken(ctx)
			if err != nil {
				return echoCtx.JSON(http.StatusUnauthorized, errors.Wrap(err, "failed to get auth token"))
			}

			authCookie := &http.Cookie{
				Name:     "session_id",
				Value:    newAuthToken.Value,
				HttpOnly: true,
				Expires:  time.Unix(newAuthToken.Expires, 0),
				SameSite: http.SameSiteLaxMode,
				Secure:   false,
				Path:     "/",
			}

			echoCtx.SetCookie(authCookie)

			newRefreshToken, err := useCase.UpdateRefresh(ctx, models.FingerPrintData{
				UserAgent: userAgent,
				UserIP:    userIP,
			})
			if err != nil {
				return echoCtx.JSON(http.StatusInternalServerError, err.Error())
			}

			cookie := &http.Cookie{
				Name:     "refresh",
				Value:    newRefreshToken.Value,
				HttpOnly: true,
				Expires:  time.Unix(newRefreshToken.Expires, 0),
				SameSite: http.SameSiteLaxMode,
				Secure:   false,
				Path:     "/",
			}

			echoCtx.SetCookie(cookie)

			return next(echoCtx)
		}
	}
}

func getCookie(echoCtx echo.Context, name string) (*http.Cookie, error) {
	cookie, err := echoCtx.Cookie(name)
	if err != nil {
		if err.Error() == "http: named cookie not present" {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get cookie")
	}

	return cookie, nil
}

func parseAuth(cookie string, authConfig userUsecase.Config) (models.AuthToken, error) {
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method")
		}
		return []byte(authConfig.AuthSignSecret), nil
	})
	if err != nil {
		return models.AuthToken{}, errors.Wrap(err, "failed to parse auth token")
	}

	if token.Valid == false {
		return models.AuthToken{}, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return models.AuthToken{}, errors.New("failed to parse token claims")
	}

	userID := claims["id"]
	expires := claims["exp"]
	if userID == nil || expires == nil {
		return models.AuthToken{}, errors.New("nil token claims")
	}

	return models.AuthToken{
		UserID:  int64(userID.(float64)),
		Expires: int64(expires.(float64)),
	}, nil
}

func getAuthToken(echoCtx echo.Context, authConfig userUsecase.Config, logger hardlogger.Logger) (models.AuthToken, error) {
	authCookie, err := getCookie(echoCtx, "session_id")
	if err != nil {
		return models.AuthToken{}, errors.Wrap(err, "failed to get auth cookie")
	}

	if authCookie == nil {
		return models.AuthToken{}, nil
	}

	authToken, err := parseAuth(authCookie.Value, authConfig)
	if err != nil {
		logger.WithError(err).Errorf("failed to parse auth token")
		return models.AuthToken{}, errors.Wrap(err, "failed to parse auth token")
	}

	return authToken, nil
}
