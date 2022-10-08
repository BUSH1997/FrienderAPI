package app

import (
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client/timepad"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/syncer"
	"github.com/BUSH1997/FrienderAPI/config"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/repository/postgres"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/usecase"
	postgreslib "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	syncer_postgres "github.com/BUSH1997/FrienderAPI/internal/pkg/syncer/repository/postgres"
	httplib "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/http"
	logger2 "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func Run() {
	configApp := config.Config{}

	err := config.LoadConfig(&configApp, "config")
	if err != nil {
		panic(err)
	}

	logger := logger2.New(os.Stdout, &logrus.JSONFormatter{}, logrus.InfoLevel)

	HTTPClient, err := httplib.NewSimpleHTTPClient(configApp.Transport.HTTP)
	if err != nil {
		panic(err)
	}

	publicEventsClient := timepad.New(HTTPClient)

	db, err := postgreslib.InitDB(configApp.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	eventRepo := postgres.New(db, logger)
	eventUsecase := usecase.New(eventRepo, logger)

	syncerRepo := syncer_postgres.New(db, logger)

	publicSyncer := syncer.New(configApp.Syncer, logger, publicEventsClient, eventUsecase, syncerRepo)

	publicSyncer.RunPublicSync()
}
