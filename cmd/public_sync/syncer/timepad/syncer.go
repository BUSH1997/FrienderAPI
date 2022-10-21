package timepad

import (
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/syncer"
)

type TimePadSyncer struct {
	syncData client.SyncData
	client   client.PublicEventsClient
}

func (ts TimePadSyncer) SyncData() client.SyncData {
	return ts.syncData
}

func (ts TimePadSyncer) Client() client.PublicEventsClient {
	return ts.client
}

type TimePadSyncData struct {
	URL string
}

func (s TimePadSyncData) GetURLs() []string {
	return []string{s.URL}
}

func (s TimePadSyncData) GetFormData() []map[string]string {
	return nil
}

func NewData(url string) client.SyncData {
	return &TimePadSyncData{
		URL: url,
	}
}

func New(syncData client.SyncData, client client.PublicEventsClient) syncer.Syncer {
	return &TimePadSyncer{
		client:   client,
		syncData: syncData,
	}
}
