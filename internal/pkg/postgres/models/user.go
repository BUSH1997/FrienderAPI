package models

type User struct {
	ID                 uint `gorm:"column:id"`
	Uid                int  `gorm:"column:uid"`
	CurrentStatus      int  `gorm:"column:current_status"`
	CreatedEventsCount int  `gorm:"column:created_events"`
	VisitedEventsCount int  `gorm:"column:visited_events"`
	IsGroup            bool `gorm:"column:is_group"`
}

func (User) TableName() string {
	return "users"
}
