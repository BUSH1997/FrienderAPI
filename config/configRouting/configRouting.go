package configRouting

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/delivery/http"
	"github.com/labstack/echo/v4"
)

type ServerConfigRouting struct {
	EventHandler *http.EventHandler
}

func (sc *ServerConfigRouting) ConfigRouting(router *echo.Echo) {
	router.POST("event/create", sc.EventHandler.CreateEvent)
}
