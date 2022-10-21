package vk

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	httplib "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/http"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type VKTransportConfig struct {
	DownloadLimitBytes int64 `mapstructure:"download_limit_bytes"`
}

type HTTPVKClient struct {
	config VKTransportConfig
	client httplib.Client
}

func New(config VKTransportConfig, client httplib.Client) client.PublicEventsClient {
	return &HTTPVKClient{
		config: config,
		client: client,
	}
}

func (c HTTPVKClient) UploadPublicEvents(ctx context.Context, data client.SyncData) ([]models.Event, error) {
	if len(data.GetURLs()) < 2 {
		return nil, errors.Wrap(
			errors.New("there must be at least 2 urls to upload events from vk"),
			"failed to perform request to vk api",
		)
	}

	getVKEventsURL := data.GetURLs()[0]
	getVKEventsDataURL := data.GetURLs()[1]

	getVKEventsFormData := data.GetFormData()[0]
	getVKEventsDataFormData := data.GetFormData()[1]

	respEvents := GetEventsResponse{
		downloadLimitBytes: c.config.DownloadLimitBytes,
	}
	err := c.client.PerformRequest(ctx, GetEventsRequestWithBody{
		GetEventsRequest: GetEventsRequest{
			RequestURL: getVKEventsURL,
		},
		FormData: getVKEventsFormData,
	}, &respEvents)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request for events to vk api")
	}

	eventIDs := make([]string, 0, len(respEvents.VKEvents.VKEventsResponse.Items))
	for _, item := range respEvents.VKEvents.VKEventsResponse.Items {
		eventIDs = append(eventIDs, strconv.Itoa(item.ID))
	}

	getVKEventsDataFormData["group_ids"] = strings.Join(eventIDs, ",")

	respEventsData := GetEventsDataResponse{
		downloadLimitBytes: c.config.DownloadLimitBytes,
	}
	err = c.client.PerformRequest(ctx, GetEventsDataRequestWithBody{
		GetEventsDataRequest: GetEventsDataRequest{
			RequestURL: getVKEventsDataURL,
		},
		FormData: getVKEventsDataFormData,
	}, &respEventsData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform request for events data to vk api")
	}

	// fmt.Println(respEventsData.VKEventsData[0].Name)

	return convertEventsToModel(respEventsData.VKEventsData.VKEventsData), nil
}

func convertEventsToModel(vkEvents []VKEventData) []models.Event {
	events := make([]models.Event, 0, len(vkEvents))
	for _, vkEvent := range vkEvents {
		event := models.Event{
			Uid:         strconv.Itoa(int(vkEvent.ID)),
			Title:       vkEvent.Name,
			StartsAt:    vkEvent.StartDate,
			IsPublic:    true,
			Description: vkEvent.Description,
			Images:      []string{vkEvent.Photo200},
		}

		if vkEvent.Photo200 == "" {
			event.Images = []string{"https://friender.hb.bizmrg.com/62f5e7ed-fb13-49e9-8af8-9ef627e697d1.jpeg"}
		}

		events = append(events, event)
	}

	return events
}
