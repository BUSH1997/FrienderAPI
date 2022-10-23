package models

type Group struct {
	ID      uint `gorm:"column:id"`
	UserId  int  `gorm:"column:user_id"`
	GroupId int  `gorm:"column:group_id"`
}

func (Group) TableName() string {
	return "groups"
}
