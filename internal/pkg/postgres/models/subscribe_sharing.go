package models

type SubscribeSharing struct {
	ID           uint `gorm:"column:id"`
	UserID       int  `gorm:"column:user_id"`
	SubscriberID int  `gorm:"column:subscriber_id"`
}

func (SubscribeSharing) TableName() string {
	return "subscribe_sharings"
}
