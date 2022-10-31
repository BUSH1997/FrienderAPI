package models

type Profile struct {
	ProfileStatus Status  `json:"profile_status,omitempty"`
	Awards        []Award `json:"awards,omitempty"`
	ActiveEvents  []Event `json:"active_events,omitempty"`
	VisitedEvents []Event `json:"visited_events,omitempty"`
}

type ChangeProfile struct {
	ProfileId int64  `json:"-"`
	Status    Status `json:"new_status_id"`
}

type UserId struct {
	Id int `json:"user_id"`
}

type Subscriptions struct {
	Groups []int64 `json:"group"`
	Users  []int64 `json:"users"`
}
