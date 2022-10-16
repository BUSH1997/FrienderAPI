package models

type UnlockedAward struct {
	ID      uint `gorm:"id"`
	UserID  int  `gorm:"user_id"`
	AwardID int  `gorm:"award_id"`
}

func (UnlockedAward) TableName() string {
	return "unlocked_awards"
}
