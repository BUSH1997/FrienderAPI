package postgres

import (
	"context"
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func (r eventRepository) Create(ctx context.Context, event models.Event) error {
	dbEvent := Event{
		Uid:         event.Uid,
		Title:       event.Title,
		Description: event.Description,
		StartsAt:    event.StartsAt,
		TimeCreated: time.Now(),
		TimeUpdated: time.Now(),
		IsGroup:     event.IsGroup,
		IsPublic:    event.IsPublic,
	}

	dbCategory := Category{}

	res := r.db.First(&dbCategory, "name = ?", event.Category)
	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to get event category")
	}

	dbEvent.Category = int(dbCategory.ID)
	dbEvent.Images = strings.Join(event.Images, ",")
	dbEvent.Geo = strconv.Itoa(int(event.GeoData.Longitude)) + "," + strconv.Itoa(int(event.GeoData.Latitude))

	dbUser := User{}

	res = r.db.Take(&dbUser, "uid = ?", event.Author)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		dbUser.Uid = event.Author
		fmt.Println("LOL1", dbUser)
		err := r.createUser(ctx, &dbUser)
		if err != nil {
			return errors.Wrap(err, "failed to create user")
		}
	}
	if err := res.Error; err != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "failed to get event category")
	}

	dbEvent.Owner = int(dbUser.ID)

	res = r.db.Create(&dbEvent)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to create event, uid %s", event.Uid)
	}

	dbEventSharings := EventSharing{}
	dbEventSharings.EventID = int(dbEvent.ID)
	dbEventSharings.UserID = int(dbUser.ID)

	res = r.db.Create(&dbEventSharings)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to create event sharing")
	}

	return nil
}

func (r eventRepository) createUser(ctx context.Context, user *User) error {
	res := r.db.Create(user)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to create user")
	}

	return nil
}
