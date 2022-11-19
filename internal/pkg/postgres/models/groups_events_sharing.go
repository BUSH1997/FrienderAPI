package models

type GroupsEventsSharing struct {
	ID            uint  `gorm:"column:id"`
	EventID       uint  `gorm:"column:event_id"`
	GroupID       uint  `gorm:"column:group_id"`
	IsAdmin       bool  `gorm:"column:is_admin"`
	IsNeedApprove bool  `gorm:"column:is_need_approve"`
	IsDeleted     bool  `gorm:"column:is_deleted"`
	UserUID       int64 `gorm:"column:user_uid"`
	IsFork        bool  `gorm:"column:is_fork"`
}

func (GroupsEventsSharing) TableName() string {
	return "groups_events_sharing"
}
