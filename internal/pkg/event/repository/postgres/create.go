package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func (r eventRepository) Create(ctx context.Context, event models.Event) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		dbEvent := db_models.Event{
			Uid:         event.Uid,
			Title:       event.Title,
			Description: event.Description,
			StartsAt:    event.StartsAt,
			TimeCreated: time.Now().Unix(),
			TimeUpdated: time.Now().Unix(),
			IsPublic:    event.IsPublic,
		}

		dbCategory := db_models.Category{}

		res := r.db.Take(&dbCategory, "name = ?", event.Category)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			dbCategory.Name = string(event.Category)
			err := r.createCategory(ctx, &dbCategory)
			if err != nil {
				return errors.Wrap(err, "failed to create category")
			}
		}
		if err := res.Error; err != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "failed to get event category")
		}

		dbEvent.Category = int(dbCategory.ID)
		dbEvent.Images = strings.Join(event.Images, ",")
		dbEvent.Geo = strconv.Itoa(int(event.GeoData.Longitude)) + "," + strconv.Itoa(int(event.GeoData.Latitude))

		dbUser := db_models.User{}

		res = r.db.Take(&dbUser, "uid = ?", event.Author)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			dbUser.Uid = event.Author
			dbUser.CurrentStatus = 1
			err := r.createUser(ctx, &dbUser)
			if err != nil {
				return errors.Wrap(err, "failed to create user")
			}
		}
		if err := res.Error; err != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.Wrap(err, "failed to get event user")
		}

		dbEvent.Owner = int(dbUser.ID)
		dbEvent.CountMembers = 1

		res = r.db.Create(&dbEvent)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to create event, uid %s", event.Uid)
		}

		dbEventSharings := db_models.EventSharing{}
		dbEventSharings.EventID = int(dbEvent.ID)
		dbEventSharings.UserID = int(dbUser.ID)

		res = r.db.Create(&dbEventSharings)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to create event sharing")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func (r eventRepository) createUser(ctx context.Context, user *db_models.User) error {
	res := r.db.Create(user)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to create user")
	}

	return nil
}

func (r eventRepository) createCategory(ctx context.Context, category *db_models.Category) error {
	res := r.db.Create(category)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to create category")
	}

	return nil
}
