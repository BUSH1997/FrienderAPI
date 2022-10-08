package main

import (
	"github.com/BUSH1997/FrienderAPI/config"
	"github.com/BUSH1997/FrienderAPI/config/configRouting"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/delivery/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/repository/postgres"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/usecase"
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

	serverRouting := configRouting.ServerConfigRouting{
		EventHandler: eventHandler,
	}
	serverRouting.ConfigRouting(router)

	router.Logger.Fatal(router.Start(":8080"))
}