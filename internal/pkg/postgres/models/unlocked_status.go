package models

type UnlockedStatus struct {
	ID       uint `gorm:"column:id"`
	UserID   int  `gorm:"column:user_id"`
	StatusID int  `gorm:"column:status_id"`
}

func (UnlockedStatus) TableName() string {
	return "unlocked_statuses"
}
