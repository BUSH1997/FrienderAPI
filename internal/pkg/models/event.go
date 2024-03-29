package models

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	_ "time"
)

const (
	SOURSE_EVENT_GROUP      = "group"
	SOURCE_EVENT_FORK_GROUP = "fork_group"
	SOURCE_EVENT_USER       = "user"
	SOURCE_EVENT_VK         = "vk_event"
)

const (
	UPDATE_ALBUM_DELETE = "delete"
	UPDATE_ALBUM_ADD    = "add"
)

type UpdateAlbumInfo struct {
	UidAlbum string `json:"uid_album"`
	UidEvent string `json:"uid_event"`
	Type     string `json:"type"`
}

type UidEventPriority struct {
	UidUser  int    `json:"-"`
	UidEvent string `json:"event_uid,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

type PriorityEvent struct {
	Priority int   `json:"priority,omitempty"`
	Event    Event `json:"event,omitempty"`
}

type GroupInfo struct {
	IsAdmin       bool  `json:"is_admin, omitempty"`
	GroupId       int64 `json:"group_id, omitempty"`
	IsNeedApprove bool  `json:"is_need_approve,omitempty"`
	CheckExist    bool  `json:"-"`
}

type Event struct {
	Uid           string    `json:"id,omitempty"`
	Title         string    `json:"title,omitempty"`
	Description   string    `json:"description,omitempty"`
	Members       []int     `json:"members"`
	Images        []string  `json:"images"`
	TimeCreated   time.Time `json:"time_created,omitempty"`
	TimeUpdated   time.Time `json:"time_update,omitempty"`
	GeoData       Geo       `json:"geo,omitempty"`
	Author        int       `json:"author,omitempty"`
	StartsAt      int64     `json:"time_start,omitempty"`
	IsGroup       bool      `json:"is_group,omitempty"`
	IsPublic      bool      `json:"is_public,omitempty"`
	IsPrivate     bool      `json:"is_private,omitempty"`
	IsActive      bool      `json:"is_active,omitempty"`
	Category      Category  `json:"category,omitempty"`
	Avatar        Avatar    `json:"avatar"`
	MembersLimit  int       `json:"members_limit"`
	GroupInfo     GroupInfo `json:"group_info,omitempty"`
	Source        string    `json:"source,omitempty"`
	Ticket        Ticket    `json:"ticket,omitempty"`
	Parent        string    `json:"parent,omitempty"`
	Albums        []string  `json:"albums"`
	CanBeReported bool      `json:"can_be_reported,omitempty"`
	BlackList     []int64   `json:"blacklist,omitempty"`
}

type Ticket struct {
	Link string `json:"link,omitempty"`
	Cost string `json:"cost,omitempty"`
}

type GetEventParams struct {
	UserID        int64       `json:"id,omitempty"`
	IsOwner       Bool        `json:"is_owner,omitempty"`
	IsActive      Bool        `json:"is_active,omitempty"`
	IsSubscriber  Bool        `json:"is_subscriber,omitempty"`
	GroupId       int64       `json:"group_id,omitempty"`
	IsNeedApprove Bool        `json:"is_need_approve,omitempty"`
	IsAdmin       Bool        `json:"is_admin,omitempty"`
	Source        string      `json:"source,omitempty"`
	City          string      `json:"city,omitempty"`
	Category      Category    `json:"category,omitempty"`
	SortMembers   string      `json:"sort_members,omitempty"`
	Search        SearchInput `json:"search,omitempty"`
	Limit         int         `json:"limit,omitempty"`
	Page          int         `json:"page,omitempty"`
	UIDs          []string    `json:"-"`
	IsPublic      Bool        `json:"-"`
}

type Category string

const (
	ART Category = "SPORT"
)

type Geo struct {
	Longitude float64 `json:"longitude,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Address   string  `json:"address,omitempty"`
}

type Avatar struct {
	AvatarUrl  string `json:"avatar_url"`
	AvatarVkId string `json:"avatar_vk_id"`
}

type SubscribeType struct {
	Id      int64
	IsGroup bool
}

func (e Event) GetEtag() string {
	s := []byte(e.Title + e.Description)

	hasher := sha256.New()
	hasher.Write(s)

	return hex.EncodeToString(hasher.Sum(nil))
}

type UnsubscribeEventInput struct {
	User int64 `json:"user,omitempty"`
}

type EventList []Event
type EventFilterFunc func(Event) bool

func (l EventList) Filter(filters ...EventFilterFunc) EventList {
	if len(filters) == 0 {
		return l
	}

	filtered := make([]Event, 0, len(l))
	for _, e := range l {
		for _, filter := range filters {
			if !filter(e) {
				continue
			}
		}

		filtered = append(filtered, e)
	}

	return filtered
}
