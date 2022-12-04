package models

import "github.com/lib/pq"

type RevindexWord struct {
	ID     uint           `gorm:"column:id"`
	Word   string         `gorm:"column:word"`
	Events pq.StringArray `gorm:"column:events;type:text[]"`
}

func (RevindexWord) TableName() string {
	return "revindex_words"
}

type RevindexEvent struct {
	ID  uint   `gorm:"column:id"`
	UID string `gorm:"column:uid"`
}

func (RevindexEvent) TableName() string {
	return "revindex_events"
}
