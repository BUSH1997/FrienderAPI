package postgres

import (
	"context"
	"fmt"
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
		Where("event_sharings.is_deleted = ?", false).
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
		Uid:          dbEvent.Uid,
		Title:        dbEvent.Title,
		Description:  dbEvent.Description,
		TimeCreated:  time.Unix(dbEvent.TimeCreated, 0),
		TimeUpdated:  time.Unix(dbEvent.TimeUpdated, 0),
		Author:       dbUser.Uid,
		StartsAt:     dbEvent.StartsAt,
		IsPublic:     dbEvent.IsPublic,
		Category:     models.Category(dbCategory.Name),
		MembersLimit: dbEvent.MembersLimit,
		Avatar: models.Avatar{
			AvatarUrl:  dbEvent.AvatarUrl,
			AvatarVkId: dbEvent.AvatarVkId,
		},
	}

	event.Members = members
	fmt.Println(strings.Split(dbEvent.Geo, ",")[0])
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
	event.IsActive = event.StartsAt > time.Now().Unix()
	images := strings.Split(dbEvent.Images, ",")
	event.Images = images

	return event, nil
}

func (r eventRepository) GetAll(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvents []db_models.Event

		query := r.db.Where("is_deleted = ?", false)

		if params.IsActive.IsDefinedTrue() {
			query = query.Where("starts_at > ?", time.Now().Unix())
		}
		if params.IsActive.IsDefinedFalse() {
			query = query.Where("starts_at < ?", time.Now().Unix())
		}

		res := query.Find(&dbEvents)
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

func (r eventRepository) GetUserEvents(ctx context.Context, user int64) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvents []db_models.Event
		res := r.db.Model(&db_models.Event{}).
			Joins("JOIN event_sharings on event_sharings.event_id = events.id").
			Joins("JOIN users on event_sharings.user_id = users.id").
			Find(&dbEvents, "users.uid = ? AND is_deleted = ?", user, false)
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

func (r eventRepository) GetUserActiveEvents(ctx context.Context, user int64) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEventSharings []db_models.EventSharing
		res := r.db.Model(&db_models.EventSharing{}).
			Joins("JOIN users on event_sharings.user_id = users.id").
			Joins("JOIN events on event_sharings.event_id = events.id").
			Where("users.uid = ?", user).
			Where("events.starts_at >= ?", time.Now().Unix()).
			Where("event_sharings.is_deleted = ?", false).
			Find(&dbEventSharings)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user event sharings")
		}

		if len(dbEventSharings) == 0 {
			return nil
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

func (r eventRepository) GetUserVisitedEvents(ctx context.Context, user int64) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEventSharings []db_models.EventSharing
		res := r.db.Model(&db_models.EventSharing{}).
			Joins("JOIN users on event_sharings.user_id = users.id").
			Joins("JOIN events on event_sharings.event_id = events.id").
			Where("users.uid = ?", user).
			Where("events.starts_at < ?", time.Now().Unix()).
			Where("event_sharings.is_deleted = ?", false).
			Find(&dbEventSharings)
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

func (r eventRepository) GetOwnerEvents(ctx context.Context, user int64) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvents []db_models.Event
		res := r.db.Model(&db_models.Event{}).
			Joins("JOIN users on events.owner_id = users.id").
			Where("users.uid = ?", user).
			Where("is_deleted = ?", false).
			Find(&dbEvents)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get owner events")
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
		return nil, errors.Wrap(err, "failed to make transaction GetOwnerEvents")
	}

	return ret, nil
}

func (r eventRepository) GetSubscriptionEvents(ctx context.Context, user int64) ([]models.Event, error) {
	var ret []models.Event

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbSubscribeSharings []db_models.SubscribeSharing
		res := r.db.Model(&db_models.SubscribeSharing{}).
			Joins("JOIN users on subscribe_sharings.subscriber_id = users.id").
			Find(&dbSubscribeSharings, "users.uid = ?", user)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user subscribe sharings")
		}

		subscribeSharingIDs := make([]int, 0, len(dbSubscribeSharings))
		for _, dbSubscribeSharing := range dbSubscribeSharings {
			subscribeSharingIDs = append(subscribeSharingIDs, dbSubscribeSharing.SubscriberID)
		}

		var dbSubscriptions []db_models.SubscribeSharing
		res = r.db.Find(&dbSubscriptions, subscribeSharingIDs)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user subscriptions")
		}

		for _, dbSubscription := range dbSubscriptions {
			subscriptionEvents, err := r.GetUserActiveEvents(ctx, int64(dbSubscription.UserID))
			if err != nil {
				return errors.Wrap(err, "failed to get subscription events")
			}

			ret = append(ret, subscriptionEvents...)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return ret, nil
}

func (r eventRepository) GetGroupEvent(ctx context.Context, group int64, isActive models.Bool) ([]models.Event, error) {
	var ret []models.Event
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbGroup db_models.Group
		res := r.db.Take(&dbGroup).Where("group_id = ?", group)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get groups events sharings")
		}

		var dbEvents []db_models.Event
		res = r.db.Model(&db_models.Event{}).
			Joins("JOIN groups_events_sharing on groups_events_sharing.event_id = events.id").
			Where("groups_events_sharing.group_id = ?", dbGroup.ID).
			Find(&dbEvents)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get events group")
		}

		ret = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			event, err := r.GetEventById(ctx, dbEvent.Uid)
			if err != nil {
				return errors.Wrapf(err, "failed to get event by id %s", dbEvent.Uid)
			}

			if !isActive.IsDefined() {
				ret = append(ret, event)
			} else if isActive.IsDefinedTrue() && event.IsActive {
				ret = append(ret, event)
			} else if isActive.IsDefinedFalse() && !event.IsActive {
				ret = append(ret, event)
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return ret, nil
}

func (r eventRepository) GetGroupAdminEvent(ctx context.Context, group int64, isAdmin models.Bool, isActive models.Bool) ([]models.Event, error) {
	var ret []models.Event
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbGroup db_models.Group
		res := r.db.Take(&dbGroup).Where("group_id = ?", group, isAdmin.Value)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get groups events sharings")
		}

		var dbEvents []db_models.Event
		res = r.db.Model(&db_models.Event{}).
			Joins("JOIN groups_events_sharing on groups_events_sharing.event_id = events.id").
			Where("groups_events_sharing.group_id = ? and groups_events_sharing.is_admin = ?", dbGroup.ID, isAdmin.Value).
			Find(&dbEvents)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get events group")
		}

		ret = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			event, err := r.GetEventById(ctx, dbEvent.Uid)
			if err != nil {
				return errors.Wrapf(err, "failed to get event by id %s", dbEvent.Uid)
			}

			if !isActive.IsDefined() {
				ret = append(ret, event)
			} else if isActive.IsDefinedTrue() && event.IsActive {
				ret = append(ret, event)
			} else if isActive.IsDefinedFalse() && !event.IsActive {
				ret = append(ret, event)
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return ret, nil
}
