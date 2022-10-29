package models

type Search struct {
	Words  []string `json:"words,omitempty"`
	Source string   `json:"source,omitempty"`
}
