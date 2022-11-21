package configMiddleware

import (
	custommiddleware "github.com/BUSH1997/FrienderAPI/cmd/main_server/middleware"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/user"
	userUsecase "github.com/BUSH1997/FrienderAPI/internal/pkg/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

var (
	allowOrigins  = []string{"*"}
	allowMethods  = []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions}
	allowHeaders  = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Csrf-Token", "X-User-Id"}
	exposeHeaders = []string{"Authorization", "X-Csrf-Token", "X-User-Id"}
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

func ConfigMiddleware(router *echo.Echo, profileRepository profile.Repository, userUseCase user.UseCase, authConfig userUsecase.Config, logger hardlogger.Logger) {
	router.Use(
		middleware.CORSWithConfig(GetCORSConfigStruct()),
		custommiddleware.UserID(logger),
		custommiddleware.CreateUser(profileRepository, logger),
		custommiddleware.RequestID(),
		custommiddleware.CheckSession(logger, userUseCase, authConfig),
	)
}
