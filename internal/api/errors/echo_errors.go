package errors

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ConfigErrorHandler(router *echo.Echo, logger hardlogger.Logger) {
	router.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Message == "missing csrf token in request header" {
				logger.WithCtx(c.Request().Context()).WithError(err).Errorf("missing csrf token in request header")
				err = c.JSON(http.StatusUnauthorized, ForbiddenError{
					Cause:  he,
					Reason: CSRFTokenEmpty,
				})
			}
			if he.Message == "invalid csrf token" {
				logger.WithCtx(c.Request().Context()).WithError(err).Errorf("invalid csrf token")
				err = c.JSON(http.StatusUnauthorized, ForbiddenError{
					Cause:  he,
					Reason: CSRFTokenInvalid,
				})
			}

		}

		router.DefaultHTTPErrorHandler(err, c)
	}
}
