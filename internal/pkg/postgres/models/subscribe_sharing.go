package models

type SubscribeSharing struct {
	Id           int
	UserId       int
	SubscriberId int
}

func (SubscribeSharing) TableName() string {
	return "subscribe_sharing"
}
