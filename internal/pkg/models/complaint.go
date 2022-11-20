package models

type Complaint struct {
	User   int64  `json:"user,omitempty"`
	Event  string `json:"event,omitempty"`
	Reason string `json:"reason,omitempty"`
}
