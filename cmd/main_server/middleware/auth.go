package middleware

import (
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

func Auth(logger hardlogger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			var userIDString string

			path := context.Path()
			// logger.Println(path)
			if strings.Contains(path, "ws/") {
				userIDString = context.QueryParam("user_id")
			} else {
				userIDString = context.Request().Header.Get("X-User-ID")
			}

			if userIDString == "" {
				logger.WithError(errors.New("empty X-User-ID")).Errorf("failed to get user header")
				return context.NoContent(http.StatusUnauthorized)
			}

			userID, err := strconv.ParseInt(userIDString, 10, 32)
			if err != nil {
				logger.WithError(err).Errorf("failed to parse user id")
				return context.NoContent(http.StatusBadRequest)
			}

			ctx := context.Request().Context()
			ctx = hardlogger.AddCtxFields(ctx, hardlogger.Fields{
				"user": userID,
			})

			context.SetRequest(
				context.Request().WithContext(contextlib.SetUser(ctx, userID)),
			)

			return next(context)
		}
	}
}
