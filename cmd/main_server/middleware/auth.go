package middleware

import (
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func Auth(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			userIDString := context.Request().Header.Get("X-User-ID")
			if userIDString == "" {
				logger.WithError(errors.New("empty X-User-ID")).Errorf("failed to get user header")
				return context.NoContent(http.StatusUnauthorized)
			}

			userID, err := strconv.ParseInt(userIDString, 10, 32)
			if err != nil {
				logger.WithError(err).Errorf("failed to parse user id")
				return context.NoContent(http.StatusBadRequest)
			}

			context.SetRequest(
				context.Request().WithContext(contextlib.SetUser(context.Request().Context(), userID)),
			)

			return next(context)
		}
	}
}
