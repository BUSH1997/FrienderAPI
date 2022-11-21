package models

type User struct {
	Auth       string `json:"-"`
	AuthExp    int64  `json:"-"`
	Refresh    string `json:"-"`
	RefreshExp int64  `json:"-"`
}
