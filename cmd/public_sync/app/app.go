package app

import (
	"context"
	vk_client "github.com/BUSH1997/FrienderAPI/cmd/public_sync/client/vk"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/syncer"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/syncer/vk"
	"github.com/BUSH1997/FrienderAPI/config"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/repository/postgres"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/repository/revindex"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event/usecase"
	postgreslib "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	syncer_postgres "github.com/BUSH1997/FrienderAPI/internal/pkg/syncer/repository/postgres"
	httplib "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/http"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/logger/hardlogger"
	"log"
)

func Run() {
	configApp := config.Config{}

	err := config.LoadConfig(&configApp, "config")
	if err != nil {
		panic(err)
	}

	logger, err := hardlogger.NewLogrusLogger(configApp.Logger)
	if err != nil {
		log.Fatal(err)
	}

	HTTPClient, err := httplib.NewSimpleHTTPClient(configApp.Transport.HTTP)
	if err != nil {
		panic(err)
	}

	//timepadSyncData := timepad.NewData(configApp.Syncer.Timepad.URL)
	//timepadClient := timepad_client.New(configApp.Transport.TimePad, HTTPClient)
	//timepadSyncer := timepad.New(timepadSyncData, timepadClient)

	vkEventsDataFormData := map[string]string{
		"access_token": "vk1.a.3v18zK0yJZRszF9FRAvhVhACDcDYPqZeeEkaehZ0k-qli2EIioZif1R4mI1cfQuwxH7cqLXG2JmDGHcf4AiTma5MpwGnhyZ3FBWjMbLqlbvCjRk1AbK8_7oWxO0DZBRySBUh2XDWCtXY6SVRRl4gDq07_U3IC-IdASY5nzcVTgZ7-qoib3C8fhoU-6I1U7-e",
		"fields":       "addresses ,description ,start_date, finish_date, place, crop_photo",
		"v":            "5.131",
	}

	db, err := postgreslib.InitDB(configApp.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	eventRepo := postgres.New(db, logger)
	eventRepo = revindex.New(db, logger, eventRepo, configApp.SkipList)
	eventUsecase := usecase.New(eventRepo, nil, nil, nil, nil, logger)
	categories, err := eventUsecase.GetAllCategories(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	vkClient := vk_client.New(configApp.Transport.VK, HTTPClient)

	var syncers []syncer.Syncer
	for _, category := range categories {
		if category == "Другое" {
			continue
		}
		vkEventsFormData := map[string]string{
			"access_token": "vk1.a.3v18zK0yJZRszF9FRAvhVhACDcDYPqZeeEkaehZ0k-qli2EIioZif1R4mI1cfQuwxH7cqLXG2JmDGHcf4AiTma5MpwGnhyZ3FBWjMbLqlbvCjRk1AbK8_7oWxO0DZBRySBUh2XDWCtXY6SVRRl4gDq07_U3IC-IdASY5nzcVTgZ7-qoib3C8fhoU-6I1U7-e",
			"q":            category,
			"type":         "event",
			"v":            "5.131",
			"offset":       "3",
			"count":        "100",
			"future":       "1",
			"sort":         "6",
		}
		vkSyncData := vk.NewData(
			configApp.Syncer.VK.GetEventsURL,
			configApp.Syncer.VK.GetEventsDataURL,
			vkEventsFormData,
			vkEventsDataFormData,
		)

		vkSyncer := vk.New(vkSyncData, vkClient)
		syncers = append(syncers, vkSyncer)
	}

	vkEventsFormDataUpdateSyncer := map[string]string{
		"access_token": "vk1.a.3v18zK0yJZRszF9FRAvhVhACDcDYPqZeeEkaehZ0k-qli2EIioZif1R4mI1cfQuwxH7cqLXG2JmDGHcf4AiTma5MpwGnhyZ3FBWjMbLqlbvCjRk1AbK8_7oWxO0DZBRySBUh2XDWCtXY6SVRRl4gDq07_U3IC-IdASY5nzcVTgZ7-qoib3C8fhoU-6I1U7-e",
		"v":            "5.131",
	}

	vkUpdateSyncData := vk.NewData(
		configApp.Syncer.VK.GetEventsURL,
		configApp.Syncer.VK.GetEventsDataURL,
		vkEventsFormDataUpdateSyncer,
		vkEventsDataFormData,
	)

	updateSyncer := vk.New(vkUpdateSyncData, vkClient)

	syncerRepo := syncer_postgres.New(db, logger)

	publicSyncer := syncer.New(configApp.Syncer, logger, syncers, updateSyncer, eventUsecase, syncerRepo)

	publicSyncer.RunPublicSync()
}
