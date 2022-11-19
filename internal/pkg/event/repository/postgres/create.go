package postgres

import (
	"context"
	"fmt"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	event_pkg "github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/postgres"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"strings"
	"time"
)

func (r eventRepository) Create(ctx context.Context, event models.Event) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		dbEvent := db_models.Event{
			Uid:          event.Uid,
			Title:        event.Title,
			Description:  event.Description,
			StartsAt:     event.StartsAt,
			TimeCreated:  time.Now().Unix(),
			TimeUpdated:  time.Now().Unix(),
			IsPublic:     event.IsPublic,
			IsPrivate:    event.IsPrivate,
			Source:       event.Source,
			AvatarUrl:    event.Avatar.AvatarUrl,
			AvatarVkId:   event.Avatar.AvatarVkId,
			Ticket:       fmt.Sprintf("%s;;%s", event.Ticket.Link, event.Ticket.Cost),
			MembersLimit: event.MembersLimit,
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
		dbEvent.Geo = fmt.Sprintf("%f;;%f;;%s", event.GeoData.Longitude, event.GeoData.Latitude, event.GeoData.Address)

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
			if postgres.ProcessError(err) == postgres.UniqueViolationError {
				err = errors.Transform(err, event_pkg.ErrAlreadyExists)
			}

			return errors.Wrapf(err, "failed to create event, uid %s", event.Uid)
		}

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
			if postgres.ProcessError(err) == postgres.UniqueViolationError {
				err = errors.Transform(err, event_pkg.ErrAlreadyExists)
			}

			return errors.Wrapf(err, "failed to create event sharing")
		}

		if event.Source == models.SOURSE_EVENT_GROUP || event.Source == models.SOURCE_EVENT_FORK_GROUP {
			var dbGroup db_models.Group
			res = r.db.Find(&dbGroup, "group_id = ?", event.GroupInfo.GroupId)
			isNeedApprove := true
			if event.GroupInfo.IsAdmin || event.Source == models.SOURCE_EVENT_FORK_GROUP {
				isNeedApprove = false
			}
			dbGroupsEventsSharing := db_models.GroupsEventsSharing{
				EventID:       dbEvent.ID,
				GroupID:       dbGroup.ID,
				IsAdmin:       event.GroupInfo.IsAdmin,
				IsNeedApprove: isNeedApprove,
				UserUID:       userID,
				IsFork:        true,
			}

			res = r.db.Create(&dbGroupsEventsSharing)
			if err := res.Error; err != nil {
				if postgres.ProcessError(err) == postgres.UniqueViolationError {
					err = errors.Transform(err, event_pkg.ErrAlreadyExists)
				}

				return errors.Wrapf(err, "failed to create group event sharing")
			}
		}

		if event.Source == models.SOURCE_EVENT_FORK_GROUP {
			var parentEvent db_models.Event
			res := r.db.Model(&parentEvent).
				Where("uid = ?", event.Parent).
				Find(&parentEvent)
			if err := res.Error; err != nil {
				return errors.Wrapf(err, "parent with this uid not found")
			}

			parentEvent.Forks = append(parentEvent.Forks, int64(dbEvent.ID))
			res = r.db.Model(&db_models.Event{}).
				Where("uid = ?", parentEvent.Uid).
				Update("forks", pq.Int64Array(parentEvent.Forks))
			if err := res.Error; err != nil {
				return errors.Wrapf(err, "error update forks events")
			}

			dbEventSharing = db_models.EventSharing{}
			dbEventSharing.EventID = int(parentEvent.ID)
			dbEventSharing.UserID = int(dbUser.ID)
			dbEventSharing.Priority = len(sharingsExist) + 1

			res = r.db.Create(&dbEventSharing)
			if err := res.Error; err != nil {
				if postgres.ProcessError(err) == postgres.UniqueViolationError {
					err = errors.Transform(err, event_pkg.ErrAlreadyExists)
				}

				return errors.Wrapf(err, "failed to create event sharing in parent event")
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
