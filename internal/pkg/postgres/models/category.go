package models

type Category struct {
	ID   uint   `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (Category) TableName() string {
	return "categories"
}
