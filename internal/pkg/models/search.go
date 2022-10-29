package models

type Search struct {
	Words   []string `json:"words,omitempty"`
	Sources []string `json:"sources,omitempty"`
}
