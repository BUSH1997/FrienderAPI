package models

type UnlockedStatus struct {
	ID       uint `gorm:"id"`
	UserID   int  `gorm:"user_id"`
	StatusID int  `gorm:"status_id"`
}

func (UnlockedStatus) TableName() string {
	return "unlocked_statuses"
}
