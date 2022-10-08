package models

import "time"

type Event struct {
	ID          uint      `gorm:"column:id"`
	Uid         string    `gorm:"column:uid"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`
	Category    int       `gorm:"column:category_id"`
	Images      string    `gorm:"column:images"`
	StartsAt    int64     `gorm:"column:starts_at"`
	TimeCreated time.Time `gorm:"column:time_created"`
	TimeUpdated time.Time `gorm:"column:time_updated"`
	Geo         string    `gorm:"column:geo"`
	Owner       int       `gorm:"column:owner_id"`
	IsGroup     bool      `gorm:"column:is_group"`
	IsPublic    bool      `gorm:"column:is_public"`
}

func (Event) TableName() string {
	return "events"
}
