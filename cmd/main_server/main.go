package main

import (
	"github.com/BUSH1997/FrienderAPI/config"
	"github.com/BUSH1997/FrienderAPI/config/configMiddleware"
	"github.com/BUSH1997/FrienderAPI/config/configRouting"
	"github.com/BUSH1997/FrienderAPI/config/configValidator"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/blacklist/text_blacklist"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/chat"
	chatHandler "github.com/BUSH1997/FrienderAPI/internal/pkg/chat/delivery/http"
	chatPostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/chat/repository/postgres"
	chatUsecase "github.com/BUSH1997/FrienderAPI/internal/pkg/chat/usecase"
	complaintHandler "github.com/BUSH1997/FrienderAPI/internal/pkg/complaint/delivery/http"
	complaintPostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/complaint/repository/postgres"
	complaintUsecase "github.com/BUSH1997/FrienderAPI/internal/pkg/complaint/usecase"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/delivery/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/repository/postgres"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/repository/revindex"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/usecase"
	groupHandler "github.com/BUSH1997/FrienderAPI/internal/pkg/group/delivery/http"
	groupPostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/group/repository/postgres"
	groupUseCase "github.com/BUSH1997/FrienderAPI/internal/pkg/group/usecase"
	image "github.com/BUSH1997/FrienderAPI/internal/pkg/image/delivery/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image/repository/s3"
	imageUseCase "github.com/BUSH1997/FrienderAPI/internal/pkg/image/usecase"
	postgreslib "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	profileHandler "github.com/BUSH1997/FrienderAPI/internal/pkg/profile/delivery/http"
	profilePostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/profile/repository/postgres"
	profileUseCase "github.com/BUSH1997/FrienderAPI/internal/pkg/profile/usecase"
	searchPostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/search/repository/postgres"
	httplib "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/vk_api"
	"github.com/labstack/echo/v4"
	"log"
)

var (
	router = echo.New()
)

func main() {
	configApp := config.Config{}
	err := config.LoadConfig(&configApp, "config")
	if err != nil {
		panic(err)
	}

	vk := vk_api.VKApi{
		AccessToken: configApp.Vk.AccessToken,
		GroupId:     configApp.Vk.GroupId,
		AlbumId:     configApp.Vk.AlbumId,
		Version:     configApp.Vk.Version,
	}

	db, err := postgreslib.InitDB(configApp.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	blackLister, err := text_blacklist.New(configApp.BlackList)
	if err != nil {
		log.Fatal(err)
	}

	// logger := logger2.New(os.Stdout, &logrus.JSONFormatter{}, logrus.InfoLevel)
	logger, err := hardlogger.NewLogrusLogger(configApp.Logger)
	if err != nil {
		log.Fatal(err)
	}

	profileRepo := profilePostgres.New(db, logger)
	eventRepo := postgres.New(db, logger)
	eventRepo = revindex.New(db, logger, eventRepo, configApp.SkipList)

	searchRepo := searchPostgres.New(db, logger)

	eventUsecase := usecase.New(eventRepo, profileRepo, searchRepo, blackLister, configApp.SkipList, logger)
	eventHandler := http.NewEventHandler(eventUsecase, logger)

	imageRepo := s3.New(logger)
	imageUseCase := imageUseCase.New(imageRepo, eventRepo, logger, vk)
	imageHandler := image.NewImageHandler(imageUseCase, logger)

	HTTPClient, err := httplib.NewSimpleHTTPClient(configApp.Transport.HTTP)
	if err != nil {
		panic(err)
	}
	profileUseCase := profileUseCase.New(profileRepo, eventRepo, logger, HTTPClient)
	profileHandler := profileHandler.NewProfileHandler(profileUseCase, logger)

	groupRepo := groupPostgres.New(db, logger)
	groupUseCase := groupUseCase.New(logger, groupRepo, profileRepo)
	groupHandler := groupHandler.New(logger, groupUseCase)

	messenger := chat.NewMessenger()

	chatRepo := chatPostgres.New(db, logger)
	chatUseCase := chatUsecase.New(chatRepo, logger)
	chatHandler := chatHandler.NewChatHandler(chatUseCase, messenger, logger)

	complaintRepo := complaintPostgres.New(db, logger)
	complaintUseCase := complaintUsecase.New(complaintRepo, logger)
	complaintHandler := complaintHandler.NewComplaintHandler(complaintUseCase, logger)

	serverRouting := configRouting.ServerConfigRouting{
		EventHandler:     eventHandler,
		ImageHandler:     imageHandler,
		ProfileHandler:   profileHandler,
		GroupHandler:     groupHandler,
		ChatHandler:      chatHandler,
		ComplaintHandler: complaintHandler,
	}

	configValidator.ConfigValidator(router)
	configMiddleware.ConfigMiddleware(router, profileRepo, logger)
	serverRouting.ConfigRouting(router)

	router.Logger.Fatal(router.Start("localhost:8090"))
}
