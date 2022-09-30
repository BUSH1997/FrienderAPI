package postgres

import "time"

type Event struct {
	ID       uint      `gorm:"id"`
	Uid      int       `gorm:"uid"`
	Title    string    `gorm:"title"`
	StartsAt time.Time `gorm:"starts_at"`
	IsPublic bool      `gorm:"is_public"`
}

func (Event) TableName() string {
	return "events"
}
