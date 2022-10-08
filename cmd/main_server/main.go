package main

import (
	"github.com/BUSH1997/FrienderAPI/config"
	"github.com/BUSH1997/FrienderAPI/config/configMiddleware"
	"github.com/BUSH1997/FrienderAPI/config/configRouting"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/delivery/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/repository/postgres"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/usecase"
	image "github.com/BUSH1997/FrienderAPI/internal/pkg/image/delivery/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/image/repository/s3"
	imageUseCase "github.com/BUSH1997/FrienderAPI/internal/pkg/image/usecase"
	postgreslib "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	logger2 "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger"
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

	db, err := postgreslib.InitDB(configApp.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	logger := logger2.New(os.Stdout, &logrus.JSONFormatter{}, logrus.InfoLevel)

	eventRepo := postgres.New(db, logger)
	eventUsecase := usecase.New(eventRepo, logger)
	eventHandler := http.NewEventHandler(eventUsecase)

	imageRepo := s3.New(logger)
	imageUseCase := imageUseCase.New(imageRepo, eventRepo, logger)
	imageHandler := image.NewImageHandler(imageUseCase)

	serverRouting := configRouting.ServerConfigRouting{
		EventHandler: eventHandler,
		ImageHandler: imageHandler,
	}
	configMiddleware.ConfigMiddleware(router)
	serverRouting.ConfigRouting(router)

	router.Logger.Fatal(router.Start("localhost:8090"))
}
