package vk

import (
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/syncer"
)

type VKSyncer struct {
	syncData client.SyncData
	client   client.PublicEventsClient
}

func (ts VKSyncer) SyncData() client.SyncData {
	return ts.syncData
}

func (ts VKSyncer) Client() client.PublicEventsClient {
	return ts.client
}

type VKSyncData struct {
	GetEventsURL          string
	GetEventsDataURL      string
	GetEventsFormData     map[string]string
	GetEventsDataFormData map[string]string
}

func (s VKSyncData) GetURLs() []string {
	return []string{s.GetEventsURL, s.GetEventsDataURL}
}

func (s VKSyncData) GetFormData() []map[string]string {
	return []map[string]string{s.GetEventsFormData, s.GetEventsDataFormData}
}

func NewData(
	getEventsURL string,
	getEventsDataURL string,
	getEventsFormData map[string]string,
	getEventsDataFormData map[string]string,
) client.SyncData {
	return &VKSyncData{
		GetEventsURL:          getEventsURL,
		GetEventsDataURL:      getEventsDataURL,
		GetEventsFormData:     getEventsFormData,
		GetEventsDataFormData: getEventsDataFormData,
	}
}

func New(syncData client.SyncData, client client.PublicEventsClient) syncer.Syncer {
	return &VKSyncer{
		client:   client,
		syncData: syncData,
	}
}
