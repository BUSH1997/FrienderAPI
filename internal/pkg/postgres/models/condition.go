package models

type Condition struct {
	ID                 uint `gorm:"id"`
	CreatedEventsCount int  `gorm:"created_events"`
	VisitedEventsCount int  `gorm:"visited_events"`
}

func (Condition) TableName() string {
	return "conditions"
}
