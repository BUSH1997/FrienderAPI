package models

import "time"

type Syncer struct {
	UpdatedAt time.Time `gorm:"updated_at"`
}

func (Syncer) TableName() string {
	return "syncer"
}
