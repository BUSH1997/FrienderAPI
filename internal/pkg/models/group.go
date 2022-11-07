package models

type GroupInput struct {
	GroupId         int  `json:"group_id,omitempty"`
	UserId          int  `json:"user_id,omitempty"`
	AllowUserEvents bool `json:"allow_user_events,omitempty"`
}

type Group struct {
	GroupId         int  `json:"group_id"`
	UserId          int  `json:"user_id"`
	AllowUserEvents bool `json:"allow_user_events"`
}

type ApproveEvent struct {
	GroupId  int    `json:"group_id"`
	Approve  bool   `json:"approve"`
	EventUid string `json:"event_uid"`
}
