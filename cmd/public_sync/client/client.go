package client

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type PublicEventsClient interface {
	UploadPublicEvents(ctx context.Context, syncData SyncData) ([]models.Event, error)
	GetCountPublicEventsWithSyncData(ctx context.Context, data SyncData) (int, error)
}

type SyncData interface {
	GetURLs() []string
	GetFormData() []map[string]string
}
