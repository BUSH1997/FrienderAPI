package models

type Profile struct {
	ProfileStatus Status  `json:"profile_status,omitempty"`
	Awards        []Award `json:"awards,omitempty"`
	ActiveEvents  []Event `json:"active_events,omitempty"`
	VisitedEvents []Event `json:"visited_events,omitempty"`
}

type ChangeProfile struct {
	ProfileId string `json:"-"`
	Status    Status `json:"new_status_id"`
}
