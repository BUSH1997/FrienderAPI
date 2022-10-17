package main

import (
	"github.com/BUSH1997/FrienderAPI/config"
	"github.com/BUSH1997/FrienderAPI/config/configMiddleware"
	"github.com/BUSH1997/FrienderAPI/config/configRouting"
	awardPostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/award/repository/postgres"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/delivery/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/repository/postgres"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/usecase"
	image "github.com/BUSH1997/FrienderAPI/internal/pkg/image/delivery/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image/repository/s3"
	imageUseCase "github.com/BUSH1997/FrienderAPI/internal/pkg/image/usecase"
	postgreslib "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	profileHandler "github.com/BUSH1997/FrienderAPI/internal/pkg/profile/delivery/http"
	profilePostgres "github.com/BUSH1997/FrienderAPI/internal/pkg/profile/repository/postgres"
	profileUseCase "github.com/BUSH1997/FrienderAPI/internal/pkg/profile/usecase"
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
	log.Println(configApp.Vk.AccessToken)
	db, err := postgreslib.InitDB(configApp.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	logger := logger2.New(os.Stdout, &logrus.JSONFormatter{}, logrus.InfoLevel)

	eventRepo := postgres.New(db, logger)
	eventUsecase := usecase.New(eventRepo, logger)
	eventHandler := http.NewEventHandler(eventUsecase, logger)

	imageRepo := s3.New(logger)
	imageUseCase := imageUseCase.New(imageRepo, eventRepo, logger, vk)
	imageHandler := image.NewImageHandler(imageUseCase)

	awardRepo := awardPostgres.New(db, logger)
	statusRepo := statusPostgres.New(db, logger)

	profileRepo := profilePostgres.New(db, logger)
	profileUseCase := profileUseCase.New(profileRepo, eventRepo, awardRepo, statusRepo, logger)
	profileHandler := profileHandler.NewProfileHandler(profileUseCase, logger)

	serverRouting := configRouting.ServerConfigRouting{
		EventHandler:   eventHandler,
		ImageHandler:   imageHandler,
		ProfileHandler: profileHandler,
	}
	configMiddleware.ConfigMiddleware(router, profileRepo, logger)
	serverRouting.ConfigRouting(router)

	router.Logger.Fatal(router.Start("localhost:8090"))
}
