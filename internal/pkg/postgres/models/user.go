package models

type User struct {
	ID                 uint `gorm:"id"`
	Uid                int  `gorm:"uid"`
	CurrentStatus      int  `gorm:"current_status"`
	CreatedEventsCount int  `gorm:"created_events"`
	VisitedEventsCount int  `gorm:"visited_events"`
}

func (User) TableName() string {
	return "users"
}
