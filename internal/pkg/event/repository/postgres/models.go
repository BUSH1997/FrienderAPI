package postgres

import "time"

type Event struct {
	ID       int       `gorm:"id"`
	Uid      string    `gorm:"uid"`
	Title    string    `gorm:"title"`
	StartsAt time.Time `gorm:"starts_at"`
	IsPublic bool      `gorm:"is_public"`
}

func (Event) TableName() string {
	return "events"
}
