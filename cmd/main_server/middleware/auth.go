package middleware

import (
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		userIDString := context.Request().Header.Get("X-User-ID")
		if userIDString == "" {
			return context.NoContent(http.StatusUnauthorized)
		}

		userID, err := strconv.ParseInt(userIDString, 10, 32)
		if err != nil {
			return context.NoContent(http.StatusBadRequest)
		}

		context.SetRequest(
			context.Request().WithContext(contextlib.SetUser(context.Request().Context(), userID)),
		)

		return next(context)
	}
}
