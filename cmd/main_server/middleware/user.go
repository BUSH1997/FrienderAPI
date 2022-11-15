package middleware

import (
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateUser(profileRepository profile.Repository, logger hardlogger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			ctx := context.Request().Context()
			userID := contextlib.GetUser(ctx)
			userExists, err := profileRepository.CheckUserExists(ctx, userID)
			if err != nil {
				return context.NoContent(http.StatusBadRequest)
			}

			if userExists {
				return next(context)
			}

			err = profileRepository.Create(ctx, userID, false)
			if err != nil {
				logger.WithError(err).Errorf("failed to create user")
				return context.JSON(http.StatusBadRequest, err)
			}

			return next(context)
		}
	}
}
