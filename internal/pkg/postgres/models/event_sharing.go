package models

type EventSharing struct {
	ID        uint `gorm:"column:id"`
	EventID   int  `gorm:"column:event_id"`
	UserID    int  `gorm:"column:user_id"`
	Priority  int  `gorm:"column:priority"`
	IsDeleted bool `gorm:"column:is_deleted"`
}

func (EventSharing) TableName() string {
	return "event_sharings"
}
