package configRouting

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/delivery/http"
	group "github.com/BUSH1997/FrienderAPI/internal/pkg/group/delivery/http"
	image "github.com/BUSH1997/FrienderAPI/internal/pkg/image/delivery/http"
	profileHandler "github.com/BUSH1997/FrienderAPI/internal/pkg/profile/delivery/http"
	"github.com/labstack/echo/v4"
)

type ServerConfigRouting struct {
	EventHandler   *http.EventHandler
	ImageHandler   *image.ImageHandler
	ProfileHandler *profileHandler.ProfileHandler
	GroupHandler   *group.GroupHandler
}

func (sc *ServerConfigRouting) ConfigRouting(router *echo.Echo) {
	router.POST("event/create", sc.EventHandler.Create)
	router.GET("event/get/:id", sc.EventHandler.GetOneEvent)
	router.GET("events", sc.EventHandler.Get)
	router.POST("image/upload", sc.ImageHandler.UploadImage)
	router.GET("profile/:id", sc.ProfileHandler.GetOneProfile)
	router.GET("profile/statuses", sc.ProfileHandler.GetAllStatusesUser)
	router.PUT("profile/change", sc.ProfileHandler.ChangeProfile)
	router.PUT("profile/events/priority", sc.ProfileHandler.ChangePriorityEvent)
	router.PUT("event/:id/subscribe", sc.EventHandler.SubscribeEvent)
	router.PUT("event/:id/unsubscribe", sc.EventHandler.UnsubscribeEvent)
	router.PUT("event/:id/delete", sc.EventHandler.DeleteEvent)
	router.PUT("event/change", sc.EventHandler.ChangeEvent)
	router.GET("categories", sc.EventHandler.GetAllCategory)
	router.POST("group/create", sc.GroupHandler.CreateGroup)
	router.GET("group", sc.GroupHandler.GetAdministeredGroup)
}
