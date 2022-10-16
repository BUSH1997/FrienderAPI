package models

type Condition struct {
	ID                 uint `gorm:"column:id"`
	CreatedEventsCount int  `gorm:"column:created_events"`
	VisitedEventsCount int  `gorm:"column:visited_events"`
}

func (Condition) TableName() string {
	return "conditions"
}
