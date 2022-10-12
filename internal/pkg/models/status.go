package models

type Status struct {
	Id       int    `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	IsLocked bool   `json:"is_locked,omitempty"`
}
