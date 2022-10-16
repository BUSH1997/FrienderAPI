package models

type SubscribeSharing struct {
	ID           uint `gorm:"id"`
	UserID       int  `gorm:"user_id"`
	SubscriberID int  `gorm:"subscriber_id"`
}

func (SubscribeSharing) TableName() string {
	return "subscribe_sharings"
}
