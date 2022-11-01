package models

type Group struct {
	GroupId         int  `json:"group_id,omitempty"`
	UserId          int  `json:"user_id,omitempty"`
	AllowUserEvents bool `json:"allow_user_events,omitempty"`
}
