package configMiddleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

var (
	allowOrigins  = []string{"*"}
	allowMethods  = []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions}
	allowHeaders  = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Csrf-Token", "X-User-ID"}
	exposeHeaders = []string{"Authorization", "X-Csrf-Token"}
)

func GetCORSConfigStruct() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    exposeHeaders,
		AllowCredentials: true,
	}
}

func ConfigMiddleware(router *echo.Echo) {
	router.Use(
		middleware.CORSWithConfig(GetCORSConfigStruct()),
	)
}
