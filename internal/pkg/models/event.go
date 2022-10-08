package models

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
	_ "time"
)

type Event struct {
	Uid         string    `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Members     []int     `json:"members,omitempty"`
	Images      []string  `json:"images,omitempty"`
	TimeCreated time.Time `json:"time_created,omitempty"`
	TimeUpdated time.Time `json:"time_update,omitempty"`
	GeoData     Geo       `json:"geo,omitempty"`
	Author      int       `json:"author,omitempty"`
	StartsAt    int64     `json:"time_start,omitempty"`
	IsGroup     bool      `json:"is_group,omitempty"`
	IsPublic    bool      `json:"is_public,omitempty"`
	Category    Category  `json:"category,omitempty"`
}

type Category string

const (
	ART Category = "SPORT"
)

type Geo struct {
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
}

func (e Event) GetEtag() string {
	s := []byte(e.Uid + e.Title + strconv.Itoa(int(e.StartsAt)))

	hasher := sha256.New()
	hasher.Write(s)

	return hex.EncodeToString(hasher.Sum(nil))
}
