package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	_ "time"
)

type UidEventPriority struct {
	UidUser  int    `json:"-"`
	UidEvent string `json:"event_uid,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

type PriorityEvent struct {
	Priority int   `json:"priority,omitempty"`
	Event    Event `json:"event,omitempty"`
}

type Event struct {
	Uid          string    `json:"id,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
	Members      []int     `json:"members,omitempty"`
	Images       []string  `json:"images,omitempty"`
	TimeCreated  time.Time `json:"time_created,omitempty"`
	TimeUpdated  time.Time `json:"time_update,omitempty"`
	GeoData      Geo       `json:"geo,omitempty"`
	Author       int       `json:"author,omitempty"`
	StartsAt     int64     `json:"time_start,omitempty"`
	IsGroup      bool      `json:"is_group,omitempty"`
	IsPublic     bool      `json:"is_public,omitempty"`
	IsPrivate    bool      `json:"is_private,omitempty"`
	IsActive     bool      `json:"is_active, omitempty"`
	Category     Category  `json:"category,omitempty"`
	Avatar       Avatar    `json:"avatar"`
	MembersLimit int       `json:"members_limit,omitempty"`
}

type GetEventParams struct {
	UserID       int64
	IsOwner      Bool
	IsActive     Bool
	IsSubscriber Bool
}

type Category string

const (
	ART Category = "SPORT"
)

type Geo struct {
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
}

type Avatar struct {
	AvatarUrl  string `json:"avatar_url"`
	AvatarVkId string `json:"avatar_vk_id"`
}

func (e Event) GetEtag() string {
	s := []byte(e.Title)

	hasher := sha256.New()
	hasher.Write(s)

	return hex.EncodeToString(hasher.Sum(nil))
}
