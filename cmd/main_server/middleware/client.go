package middleware

import (
	api_errors "github.com/BUSH1997/FrienderAPI/internal/api/errors"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ClientID(clients map[string]bool, logger hardlogger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			clientID := context.Request().Header.Get("X-Client-ID")
			if !clients[clientID] {
				err := errors.Typed("invalid_client_id", "client id is invalid")

				logger.WithCtx(context.Request().Context()).WithError(err).Errorf("client id is invalid")
				return context.JSON(http.StatusUnauthorized, api_errors.ForbiddenError{
					Cause:  err,
					Reason: api_errors.NoAccess,
				})
			}

			ctx := context.Request().Context()
			ctx = hardlogger.AddCtxFields(ctx, hardlogger.Fields{
				"client_id": clientID,
			})

			context.SetRequest(context.Request().WithContext(contextlib.SetRequestID(ctx, clientID)))

			return next(context)
		}
	}
}
