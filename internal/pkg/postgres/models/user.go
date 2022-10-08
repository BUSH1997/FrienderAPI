package models

type User struct {
	ID  uint `gorm:"id"`
	Uid int  `gorm:"uid"`
}

func (User) TableName() string {
	return "users"
}
