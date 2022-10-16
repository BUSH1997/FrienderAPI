package models

type Award struct {
	ID          uint   `gorm:"id"`
	Image       string `gorm:"image"`
	Name        string `gorm:"name"`
	Description string `gorm:"description"`
	ConditionID int    `gorm:"condition_id"`
}

func (Award) TableName() string {
	return "awards"
}
