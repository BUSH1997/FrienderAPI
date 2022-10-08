package postgres

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

type EventSharing struct {
	ID      uint `gorm:"id"`
	EventID int  `gorm:"event_id"`
	UserID  int  `gorm:"user_id"`
}

func (EventSharing) TableName() string {
	return "event_sharings"
}

type User struct {
	ID  uint `gorm:"id"`
	Uid int  `gorm:"uid"`
}

func (User) TableName() string {
	return "users"
}

type Category struct {
	ID   uint   `gorm:"id"`
	Name string `gorm:"name"`
}

func (Category) TableName() string {
	return "categories"
}
