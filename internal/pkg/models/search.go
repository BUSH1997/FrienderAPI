package models

type SearchInput struct {
	Enabled    bool   `json:"enabled,omitempty"`
	SearchData Search `json:"data,omitempty"`
}

type Search struct {
	Words   []string `json:"words,omitempty"`
	Sources []string `json:"sources,omitempty"`
}
