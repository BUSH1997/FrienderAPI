package models

import "time"

type Syncer struct {
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Syncer) TableName() string {
	return "syncer"
}
