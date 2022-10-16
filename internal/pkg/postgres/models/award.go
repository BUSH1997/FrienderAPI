package models

type Award struct {
	ID          uint   `gorm:"column:id"`
	Image       string `gorm:"column:image"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	ConditionID int    `gorm:"column:condition_id"`
}

func (Award) TableName() string {
	return "awards"
}
