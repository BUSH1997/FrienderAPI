package configRouting

import (
	chat "github.com/BUSH1997/FrienderAPI/internal/pkg/chat/delivery/http"
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
	ChatHandler    *chat.ChatHandler
}

func (sc *ServerConfigRouting) ConfigRouting(router *echo.Echo) {
	router.POST("event/create", sc.EventHandler.Create)
	router.GET("event/get/:id", sc.EventHandler.GetOneEvent)
	router.POST("events", sc.EventHandler.Get)
	router.POST("image/upload", sc.ImageHandler.UploadImage)
	router.GET("profile/friends", sc.ProfileHandler.GetFriends)
	router.GET("profile/:id", sc.ProfileHandler.GetOneProfile)
	router.GET("profile/statuses", sc.ProfileHandler.GetAllStatusesUser)
	router.PUT("profile/change", sc.ProfileHandler.ChangeProfile)
	router.POST("profile/subscribe", sc.ProfileHandler.Subscribe)
	router.POST("profile/unsubscribe", sc.ProfileHandler.UnSubscribe)
	router.GET("profile/get", sc.ProfileHandler.GetSubscribe)
	router.PUT("profile/events/priority", sc.ProfileHandler.ChangePriorityEvent)
	router.PUT("event/:id/subscribe", sc.EventHandler.SubscribeEvent)
	router.PUT("event/:id/unsubscribe", sc.EventHandler.UnsubscribeEvent)
	router.PUT("event/:id/delete", sc.EventHandler.DeleteEvent)
	router.PUT("event/change", sc.EventHandler.ChangeEvent)
	router.GET("categories", sc.EventHandler.GetAllCategory)
	router.POST("group/create", sc.GroupHandler.CreateGroup)
	router.GET("group", sc.GroupHandler.GetAdministeredGroup)
	router.GET("group/admin/check", sc.GroupHandler.IsAdmin)
	router.GET("group/get", sc.GroupHandler.Get)
	router.PUT("group/update", sc.GroupHandler.Update)
	router.GET("ws/messenger/:id", sc.ChatHandler.ProcessMessage)
	router.GET("messages", sc.ChatHandler.GetMessages)
	router.GET("chats", sc.ChatHandler.GetChats)
	router.GET("cities", sc.ProfileHandler.GetAllCities)
	router.POST("event/:id/photo/upload", sc.EventHandler.UploadPhotos)
	router.POST("event/:id/photo/delete", sc.EventHandler.DeletePhotos)
}
