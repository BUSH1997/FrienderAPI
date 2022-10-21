package app

import (
	vk_client "github.com/BUSH1997/FrienderAPI/cmd/public_sync/client/vk"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/syncer"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/syncer/vk"
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

	//timepadSyncData := timepad.NewData(configApp.Syncer.Timepad.URL)
	//timepadClient := timepad_client.New(configApp.Transport.TimePad, HTTPClient)
	//timepadSyncer := timepad.New(timepadSyncData, timepadClient)

	vkEventsFormData := map[string]string{
		"access_token": "vk1.a.3v18zK0yJZRszF9FRAvhVhACDcDYPqZeeEkaehZ0k-qli2EIioZif1R4mI1cfQuwxH7cqLXG2JmDGHcf4AiTma5MpwGnhyZ3FBWjMbLqlbvCjRk1AbK8_7oWxO0DZBRySBUh2XDWCtXY6SVRRl4gDq07_U3IC-IdASY5nzcVTgZ7-qoib3C8fhoU-6I1U7-e",
		"q":            "кино",
		"type":         "event",
		"v":            "5.131",
		"offset":       "3",
	}

	vkEventsDataFormData := map[string]string{
		"access_token": "vk1.a.3v18zK0yJZRszF9FRAvhVhACDcDYPqZeeEkaehZ0k-qli2EIioZif1R4mI1cfQuwxH7cqLXG2JmDGHcf4AiTma5MpwGnhyZ3FBWjMbLqlbvCjRk1AbK8_7oWxO0DZBRySBUh2XDWCtXY6SVRRl4gDq07_U3IC-IdASY5nzcVTgZ7-qoib3C8fhoU-6I1U7-e",
		"fields":       "addresses ,description ,start_date, finish_date",
		"v":            "5.131",
	}

	vkSyncData := vk.NewData(
		configApp.Syncer.VK.GetEventsURL,
		configApp.Syncer.VK.GetEventsDataURL,
		vkEventsFormData,
		vkEventsDataFormData,
	)
	vkClient := vk_client.New(configApp.Transport.VK, HTTPClient)
	vkSyncer := vk.New(vkSyncData, vkClient)

	db, err := postgreslib.InitDB(configApp.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	syncers := []syncer.Syncer{vkSyncer}

	eventRepo := postgres.New(db, logger)
	eventUsecase := usecase.New(eventRepo, nil, logger)

	syncerRepo := syncer_postgres.New(db, logger)

	publicSyncer := syncer.New(configApp.Syncer, logger, syncers, eventUsecase, syncerRepo)

	publicSyncer.RunPublicSync()
}
