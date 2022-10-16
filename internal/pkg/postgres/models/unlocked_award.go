package models

type UnlockedAward struct {
	ID      uint `gorm:"column:id"`
	UserID  int  `gorm:"column:user_id"`
	AwardID int  `gorm:"column:award_id"`
}

func (UnlockedAward) TableName() string {
	return "unlocked_awards"
}
