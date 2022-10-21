package postgres

import (
	"context"
	"fmt"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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
			IsPrivate:   event.IsPrivate,
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
		dbEvent.Geo = fmt.Sprintf("%f", event.GeoData.Longitude) + "," + fmt.Sprintf("%f", event.GeoData.Latitude)

		dbUser := db_models.User{}

		userID := contextlib.GetUser(ctx)

		res = r.db.Take(&dbUser, "uid = ?", userID)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get event user")
		}

		res = r.db.Model(&db_models.User{}).Where("id = ?", dbUser.ID).
			Update("created_events", dbUser.CreatedEventsCount+1)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to update user created events")
		}

		dbEvent.Owner = int(dbUser.ID)
		dbEvent.CountMembers = 1

		res = r.db.Create(&dbEvent)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to create event, uid %s", event.Uid)
		}

		if !event.IsPublic {
			var sharingsExist []db_models.EventSharing
			res = r.db.Model(&db_models.EventSharing{}).Find(&sharingsExist, "user_id = ?", dbUser.ID)
			if err := res.Error; err != nil {
				return errors.Wrap(err, "failed to get user event sharings")
			}

			dbEventSharing := db_models.EventSharing{}
			dbEventSharing.EventID = int(dbEvent.ID)
			dbEventSharing.UserID = int(dbUser.ID)
			dbEventSharing.Priority = len(sharingsExist) + 1

			res = r.db.Create(&dbEventSharing)
			if err := res.Error; err != nil {
				return errors.Wrapf(err, "failed to create event sharing")
			}
		} else {
			var dbGroup db_models.Group
			res = r.db.Find(&dbGroup, "group_id = ?", event.GroupInfo.GroupId)

			dbGroupsEventsSharing := db_models.GroupsEventsSharing{
				EventID: dbEvent.ID,
				GroupID: dbGroup.ID,
			}

			res = r.db.Create(&dbGroupsEventsSharing)
			if err := res.Error; err != nil {
				return errors.Wrapf(err, "failed to create group event sharing")
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
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
