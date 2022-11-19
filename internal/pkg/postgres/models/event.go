package models

import "github.com/lib/pq"

type Event struct {
	ID                   uint           `gorm:"column:id"`
	Uid                  string         `gorm:"column:uid"`
	Title                string         `gorm:"column:title"`
	Description          string         `gorm:"column:description"`
	Images               string         `gorm:"column:images"`
	AvatarUrl            string         `gorm:"column:avatar_url"`
	AvatarVkId           string         `gorm:"column:avatar_vk_id"`
	StartsAt             int64          `gorm:"column:starts_at"`
	TimeCreated          int64          `gorm:"column:time_created"`
	TimeUpdated          int64          `gorm:"column:time_updated"`
	Geo                  string         `gorm:"column:geo"`
	Category             int            `gorm:"column:category_id"`
	CountMembers         int            `gorm:"column:count_members"`
	IsPublic             bool           `gorm:"column:is_public"`
	IsPrivate            bool           `gorm:"column:is_private"`
	Owner                int            `gorm:"column:owner_id"`
	IsDeleted            bool           `gorm:"column:is_deleted"`
	Photos               string         `gorm:"column:photos"`
	MembersLimit         int            `gorm:"column:members_limit"`
	Source               string         `gorm:"column:source"`
	Ticket               string         `gorm:"column:ticket"`
	Forks                pq.Int64Array  `gorm:"column:forks;type:integer[]"`
	Albums               pq.StringArray `gorm:"column:albums; type:varchar(64)[]"`
	LastMessageCreatedAt int64          `gorm:"column:last_message_created_at"`
}

func (Event) TableName() string {
	return "events"
}
