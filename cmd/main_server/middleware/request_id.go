package middleware

import (
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
)

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			ctx := context.Request().Context()
			requestID := contextlib.GetOrGenerateRequestID(ctx)
			ctx = hardlogger.AddCtxFields(ctx, hardlogger.Fields{
				"request_id": requestID,
			})

			context.SetRequest(context.Request().WithContext(contextlib.SetRequestID(ctx, requestID)))

			return next(context)
		}
	}
}
