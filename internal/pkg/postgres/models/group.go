package models

type Group struct {
	ID      uint `gorm:"column:id"`
	UserId  int  `gorm:"user_id"`
	GroupId int  `gorm:"group_id"`
}

func (Group) TableName() string {
	return "groups"
}
