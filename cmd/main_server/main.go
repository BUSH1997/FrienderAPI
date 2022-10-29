package main

import (
	"github.com/BUSH1997/FrienderAPI/config"
	"github.com/BUSH1997/FrienderAPI/config/configMiddleware"
	"github.com/BUSH1997/FrienderAPI/config/configRouting"
	awardPostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/award/repository/postgres"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/blacklist/text_blacklist"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/chat"
	chatHandler "github.com/BUSH1997/FrienderAPI/internal/pkg/chat/delivery/http"
	chatPostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/chat/repository/postgres"
	chatUsecase "github.com/BUSH1997/FrienderAPI/internal/pkg/chat/usecase"
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
	searchHandler "github.com/BUSH1997/FrienderAPI/internal/pkg/search/delivery/http"
	searchPostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/search/repository/postgres"
	searchUsecase "github.com/BUSH1997/FrienderAPI/internal/pkg/search/usecase"
	statusPostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/status/repository/postgres"
	logger2 "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/vk_api"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"log"
	"os"
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

	logger := logger2.New(os.Stdout, &logrus.JSONFormatter{}, logrus.InfoLevel)
	logger.Println(configApp.Vk.AccessToken)
	profileRepo := profilePostgres.New(db, logger)
	eventRepo := postgres.New(db, logger)
	eventRepo = revindex.New(db, logger, eventRepo, configApp.SkipList)
	eventUsecase := usecase.New(eventRepo, profileRepo, blackLister, logger)
	eventHandler := http.NewEventHandler(eventUsecase, logger)

	imageRepo := s3.New(logger)
	imageUseCase := imageUseCase.New(imageRepo, eventRepo, logger, vk)
	imageHandler := image.NewImageHandler(imageUseCase)

	awardRepo := awardPostgres.New(db, logger)
	statusRepo := statusPostgres.New(db, logger)

	profileUseCase := profileUseCase.New(profileRepo, eventRepo, awardRepo, statusRepo, logger)
	profileHandler := profileHandler.NewProfileHandler(profileUseCase, logger)

	groupRepo := groupPostgres.New(db, logger)
	groupUseCase := groupUseCase.New(logger, groupRepo, profileRepo)
	groupHandler := groupHandler.New(logger, groupUseCase)

	messenger := chat.NewMessenger()

	chatRepo := chatPostgres.New(db, logger)
	chatUseCase := chatUsecase.New(chatRepo)
	chatHandler := chatHandler.NewChatHandler(chatUseCase, messenger, logger)

	searchRepo := searchPostgres.New(db, logger)
	searchUseCase := searchUsecase.New(searchRepo, eventRepo, configApp.SkipList, logger)
	searchHandler := searchHandler.NewSearchHandler(searchUseCase, logger)

	serverRouting := configRouting.ServerConfigRouting{
		EventHandler:   eventHandler,
		ImageHandler:   imageHandler,
		ProfileHandler: profileHandler,
		GroupHandler:   groupHandler,
		ChatHandler:    chatHandler,
		SearchHandler:  searchHandler,
	}

	configMiddleware.ConfigMiddleware(router, profileRepo, logger)
	serverRouting.ConfigRouting(router)

	router.Logger.Fatal(router.Start("localhost:8090"))
}
