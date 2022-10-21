package models

type GroupsEventsSharing struct {
	ID      uint `gorm:"column:id"`
	EventID uint `gorm:"column:event_id"`
	GroupID uint `gorm:"column:group_id"`
}

func (GroupsEventsSharing) TableName() string {
	return "groups_events_sharing"
}
