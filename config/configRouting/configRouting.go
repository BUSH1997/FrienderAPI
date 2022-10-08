package configRouting

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/delivery/http"
	image "github.com/BUSH1997/FrienderAPI/internal/pkg/image/delivery/http"
	"github.com/labstack/echo/v4"
)

type ServerConfigRouting struct {
	EventHandler *http.EventHandler
	ImageHandler *image.ImageHandler
}

func (sc *ServerConfigRouting) ConfigRouting(router *echo.Echo) {
	router.POST("event/create", sc.EventHandler.CreateEvent)
	router.GET("event/:id", sc.EventHandler.GetOneEvent)
	router.GET("events", sc.EventHandler.GetEvents)
	router.GET("events/:id", sc.EventHandler.GetEventsUser)
	router.POST("image/upload", sc.ImageHandler.UploadImage)
}
