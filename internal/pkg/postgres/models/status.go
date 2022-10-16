package models

type Status struct {
	ID          uint   `gorm:"id"`
	UID         int    `gorm:"uid"`
	Title       string `gorm:"title"`
	ConditionID int    `gorm:"condition_id"`
}

func (Status) TableName() string {
	return "statuses"
}
