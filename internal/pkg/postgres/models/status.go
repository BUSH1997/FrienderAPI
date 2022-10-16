package models

type Status struct {
	ID          uint   `gorm:"column:id"`
	UID         int    `gorm:"column:uid"`
	Title       string `gorm:"column:title"`
	ConditionID int    `gorm:"column:condition_id"`
}

func (Status) TableName() string {
	return "statuses"
}
