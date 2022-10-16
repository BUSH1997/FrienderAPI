package models

type EventSharing struct {
	ID        uint `gorm:"id"`
	EventID   int  `gorm:"event_id"`
	UserID    int  `gorm:"user_id"`
	Priority  int  `gorm:"priority"`
	IsDeleted bool `gorm:"is_deleted"`
}

func (EventSharing) TableName() string {
	return "event_sharings"
}
