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

func (r eventRepository) GetAllPublic(ctx context.Context) ([]models.Event, error) {
	var events []models.Event
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvents []db_models.Event

		res := r.db.Find(&dbEvents, "is_public = ?", true)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to check user")
		}

		events = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			event := models.Event{
				Uid:      dbEvent.Uid,
				Title:    dbEvent.Title,
				StartsAt: dbEvent.StartsAt,
				IsPublic: dbEvent.IsPublic,
			}

			events = append(events, event)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return events, nil
}

func (r eventRepository) GetEventById(ctx context.Context, id string) (models.Event, error) {
	var event models.Event

	err := r.db.Transaction(func(tx *gorm.DB) (err error) {
		event, err = r.getEventById(ctx, id)
		if err != nil {
			return errors.Wrap(err, "failed to get event by id")
		}

		return nil
	})
	if err != nil {
		return models.Event{}, errors.Wrap(err, "failed to make transaction")
	}

	return event, nil
}

func (r eventRepository) getEventById(ctx context.Context, id string) (models.Event, error) {
	var dbEvent db_models.Event
	res := r.db.Take(&dbEvent, "uid = ?", id)
	if err := res.Error; err != nil {
		return models.Event{}, errors.Wrap(err, "failed to get event by id")
	}

	var dbUser db_models.User
	res = r.db.Take(&dbUser, "id = ?", dbEvent.Owner)
	if err := res.Error; err != nil {
		return models.Event{}, errors.Wrap(err, "failed to get owner id")
	}

	var dbCategory db_models.Category
	res = r.db.Take(&dbCategory, "id = ?", dbEvent.Category)
	if err := res.Error; err != nil {
		return models.Event{}, errors.Wrap(err, "failed to get category id")
	}

	var dbEventSharings []db_models.EventSharing

	res = r.db.
		Joins("JOIN events on event_sharings.event_id = events.id").
		Where("events.uid = ?", id).
		Find(&dbEventSharings)
	if err := res.Error; err != nil {
		return models.Event{}, errors.Wrap(err, "failed to get event sharings")
	}

	memberDBIDs := make([]int, 0, len(dbEventSharings))
	for _, eventSharing := range dbEventSharings {
		memberDBIDs = append(memberDBIDs, eventSharing.UserID)
	}

	var dbMembers []db_models.User
	res = r.db.Find(&dbMembers, memberDBIDs)
	if err := res.Error; err != nil {
		return models.Event{}, errors.Wrap(err, "failed to get members")
	}

	members := make([]int, 0, len(dbEventSharings))
	for _, dbMember := range dbMembers {
		members = append(members, dbMember.Uid)
	}

	event := models.Event{
		Uid:         dbEvent.Uid,
		Title:       dbEvent.Title,
		Description: dbEvent.Description,
		TimeCreated: time.Unix(dbEvent.TimeCreated, 0),
		TimeUpdated: time.Unix(dbEvent.TimeUpdated, 0),
		Author:      dbUser.Uid,
		StartsAt:    dbEvent.StartsAt,
		IsPublic:    dbEvent.IsPublic,
		Category:    models.Category(dbCategory.Name),
	}

	event.Members = members

	longitude, err := strconv.ParseFloat(strings.Split(dbEvent.Geo, ",")[0], 32)
	if err != nil {
		return models.Event{}, errors.Wrap(err, "failed to parse longitude")
	}

	latitude, err := strconv.ParseFloat(strings.Split(dbEvent.Geo, ",")[1], 32)
	if err != nil {
		return models.Event{}, errors.Wrap(err, "failed to parse latitude")
	}

	event.GeoData = models.Geo{
		Longitude: longitude,
		Latitude:  latitude,
	}

	images := strings.Split(dbEvent.Images, ",")
	event.Images = images

	return event, nil
}

func (r eventRepository) GetAll(ctx context.Context) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvents []db_models.Event
		res := r.db.Find(&dbEvents)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get all events")
		}

		ret = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			event, err := r.GetEventById(ctx, dbEvent.Uid)
			if err != nil {
				return errors.Wrapf(err, "failed to get event by id %s", dbEvent.Uid)
			}

			ret = append(ret, event)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return ret, nil
}

func (r eventRepository) GetUserEvents(ctx context.Context, id int64) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvents []db_models.Event
		res := r.db.Model(&db_models.Event{}).
			Joins("JOIN event_sharings on event_sharings.event_id = events.id").
			Joins("JOIN users on event_sharings.user_id = users.id").
			Find(&dbEvents, "users.uid = ?", id)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user events")
		}

		ret = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			event, err := r.GetEventById(ctx, dbEvent.Uid)
			if err != nil {
				return errors.Wrapf(err, "failed to get event by id %s", dbEvent.Uid)
			}

			ret = append(ret, event)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return ret, nil
}

func (r eventRepository) GetUserActiveEvents(ctx context.Context, id int) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEventSharings []db_models.EventSharing
		res := r.db.Model(&db_models.EventSharing{}).
			Joins("JOIN users on event_sharings.user_id = users.id").
			Joins("JOIN events on event_sharings.event_id = events.id").
			Find(&dbEventSharings, "users.uid = ? AND events.starts_at >= ?", id, time.Now().Unix())
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user event sharings")
		}

		eventIDs := make([]int, 0, len(dbEventSharings))
		for _, sharing := range dbEventSharings {
			eventIDs = append(eventIDs, sharing.EventID)
		}

		var dbEvents []db_models.Event
		res = r.db.Find(&dbEvents, eventIDs)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get events")
		}

		ret = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			event, err := r.GetEventById(ctx, dbEvent.Uid)
			if err != nil {
				return errors.Wrapf(err, "failed to get event by uid %s", dbEvent.Uid)
			}

			ret = append(ret, event)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return ret, nil
}

func (r eventRepository) GetUserVisitedEvents(ctx context.Context, id int) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEventSharings []db_models.EventSharing
		res := r.db.Model(&db_models.EventSharing{}).
			Joins("JOIN users on event_sharings.user_id = users.id").
			Joins("JOIN events on event_sharings.event_id = events.id").
			Find(&dbEventSharings, "users.uid = ? AND events.starts_at < ?", id, time.Now().Unix())
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user event sharings")
		}

		eventIDs := make([]int, 0, len(dbEventSharings))
		for _, sharing := range dbEventSharings {
			eventIDs = append(eventIDs, sharing.EventID)
		}

		var dbEvents []db_models.Event
		res = r.db.Find(&dbEvents, eventIDs)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get events")
		}

		ret = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			event, err := r.GetEventById(ctx, dbEvent.Uid)
			if err != nil {
				return errors.Wrapf(err, "failed to get event by uid %s", dbEvent.Uid)
			}

			ret = append(ret, event)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return ret, nil
}

func (r eventRepository) GetPriority(ctx context.Context, user int, event string) (int, error) {
	var dbEventSharing db_models.EventSharing

	res := r.db.Model(&db_models.EventSharing{}).
		Joins("JOIN users on event_sharings.user_id = users.id").
		Joins("JOIN events on event_sharings.event_id = events.id").
		Take(&dbEventSharing, "users.uid = ? AND events.uid = ?", user, event)
	if err := res.Error; err != nil {
		return 0, errors.Wrap(err, "failed to get event sharing")
	}

	return dbEventSharing.Priority, nil
}

func (r eventRepository) GetAllCategories(ctx context.Context) ([]string, error) {
	var ret []string
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbCategories []db_models.Category
		res := r.db.Find(&dbCategories)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get all categories")
		}

		ret = make([]string, 0, len(dbCategories))
		for _, dbEvent := range dbCategories {
			currentCategory := dbEvent.Name
			ret = append(ret, currentCategory)
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction GetAllCategories")
	}

	return ret, nil
}
