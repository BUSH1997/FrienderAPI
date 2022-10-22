package models

type Message struct {
	ID          uint   `gorm:"column:id"`
	UserID      int    `gorm:"column:user_id"`
	UserUID     int64  `gorm:"column:user_uid"`
	TimeCreated int64  `gorm:"column:time_created"`
	Text        string `gorm:"column:text"`
	EventID     int    `gorm:"column:event_id"`
	EventUID    string `gorm:"column:event_uid"`
}

func (Message) TableName() string {
	return "messages"
}
