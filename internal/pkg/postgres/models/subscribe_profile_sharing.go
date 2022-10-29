package models

type SubscribeProfileSharing struct {
	ID        uint  `gorm:"column:id"`
	ProfileId int64 `gorm:"column:profile_id"`
	UserId    int64 `gorm:"column:user_id"`
}

func (SubscribeProfileSharing) TableName() string {
	return "subscribe_profile_sharing"
}
