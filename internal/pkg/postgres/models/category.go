package models

type Category struct {
	ID   uint   `gorm:"id"`
	Name string `gorm:"name"`
}

func (Category) TableName() string {
	return "categories"
}
