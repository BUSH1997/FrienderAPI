package models

type Award struct {
	Image       string `json:"image,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IsLocked    bool   `json:"is_locked,omitempty"`
}
