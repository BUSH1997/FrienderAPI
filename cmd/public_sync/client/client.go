package client

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
)

type PublicEventsClient interface {
	UploadPublicEvents(ctx context.Context, url string) ([]models.Event, error)
}
