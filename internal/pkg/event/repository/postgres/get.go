package postgres

import (
	"context"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

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

func (r eventRepository) findMembersByEventUid(ctx context.Context, uid string) ([]int, error) {
	var dbEventSharings []db_models.EventSharing
	res := r.db.
		Joins("JOIN events on event_sharings.event_id = events.id").
		Where("events.uid = ?", uid).
		Where("event_sharings.is_deleted = ?", false).
		Find(&dbEventSharings)
	if err := res.Error; err != nil {
		return []int{}, errors.Wrap(err, "failed to get event sharings")
	}

	memberDBIDs := make([]int, 0, len(dbEventSharings))
	for _, eventSharing := range dbEventSharings {
		memberDBIDs = append(memberDBIDs, eventSharing.UserID)
	}

	var dbMembers []db_models.User
	res = r.db.Find(&dbMembers, memberDBIDs)
	if err := res.Error; err != nil {
		return []int{}, errors.Wrap(err, "failed to get members")
	}

	members := make([]int, 0, len(dbEventSharings))
	for _, dbMember := range dbMembers {
		members = append(members, dbMember.Uid)
	}

	return members, nil
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

	members, err := r.findMembersByEventUid(ctx, id)
	for _, forkId := range dbEvent.Forks {
		res := r.db.Take(&dbEvent, "id = ?", forkId)
		if err := res.Error; err != nil {
			return models.Event{}, errors.Wrap(err, "failed to get event by id")
		}

		membersFork, err := r.findMembersByEventUid(ctx, dbEvent.Uid)
		if err != nil {
			return models.Event{}, err
		}
		members = append(members, membersFork...)
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
		Source: dbEvent.Source,
	}

	if strings.Contains(dbEvent.Ticket, ";;") {
		ticketData := strings.Split(dbEvent.Ticket, ";;")
		if len(ticketData) != 2 {
			return models.Event{}, errors.New("assumed two items in ticket data")
		}

		event.Ticket.Link = ticketData[0]
		event.Ticket.Cost = ticketData[1]
	}

	if dbEvent.Source == "group" {
		var groupEventSharing db_models.GroupsEventsSharing
		res := r.db.Take(&groupEventSharing, "event_id = ?", dbEvent.ID)
		if err := res.Error; err != nil {
			return models.Event{}, errors.Wrap(err, "failed to get groupEventSharing")
		}

		var group db_models.Group
		res = r.db.Take(&group, "id = ?", groupEventSharing.GroupID)
		if err := res.Error; err != nil {
			return models.Event{}, errors.Wrap(err, "failed to get group with id")
		}
		event.GroupInfo.GroupId = int64(group.GroupId)
		event.GroupInfo.IsAdmin = groupEventSharing.IsAdmin
		event.GroupInfo.IsNeedApprove = groupEventSharing.IsNeedApprove
	}

	event.Members = members

	geoData, err := getGeoData(dbEvent.Geo)
	if err != nil {
		return models.Event{}, errors.Wrap(err, "failed to get event geo data")
	}

	event.GeoData = geoData
	event.IsActive = event.StartsAt > time.Now().Unix()
	if dbEvent.Images != "" {
		images := strings.Split(dbEvent.Images, ",")
		event.Images = images
	} else {
		event.Images = make([]string, 0)
	}

	return event, nil
}

func getGeoData(geo string) (models.Geo, error) {
	geoData := strings.Split(geo, ";;")
	longitude, err := strconv.ParseFloat(geoData[0], 32)
	if err != nil {
		return models.Geo{}, errors.Wrap(err, "failed to parse longitude")
	}

	latitude, err := strconv.ParseFloat(geoData[1], 32)
	if err != nil {
		return models.Geo{}, errors.Wrap(err, "failed to parse latitude")
	}

	return models.Geo{
		Longitude: longitude,
		Latitude:  latitude,
		Address:   geoData[2],
	}, nil
}

func (r eventRepository) GetAll(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	var dbEvents []db_models.Event
	query := r.db.Model(&db_models.Event{})

	if params.IsOwner.IsDefinedTrue() {
		query = query.Joins("JOIN users on events.owner_id = users.id").
			Where("users.uid = ?", params.UserID)
	}

	if params.Source == "not_vk" {
		query = query.Where("source <> ?", "vk_event")
	}
	if params.Source == "vk_event" {
		query = query.Where("source = ?", params.Source)
	}
	if params.IsActive.IsDefinedTrue() {
		query = query.Where("starts_at > ?", time.Now().Unix())
	}
	if params.IsActive.IsDefinedFalse() {
		query = query.Where("starts_at < ?", time.Now().Unix())
	}

	if params.Category != "" {
		dbCategory, err := r.getCategory(string(params.Category))
		if err != nil {
			return nil, errors.Wrap(err, "failed to get category")
		}

		query = query.Where("category_id = ?", dbCategory.ID)
	}

	if params.City != "" {
		query = query.Where("geo LIKE ?", "%"+params.City+"%")
	}

	if len(params.UIDs) > 0 {
		query = query.Where("uid in ?", params.UIDs)
	}

	if params.IsPublic.IsDefinedTrue() {
		query = query.Where("is_public = ?", true)
	}

	query = query.Where("is_deleted = ?", false)

	query = query.Offset(params.Page * params.Limit)

	if params.Limit != 0 {
		query = query.Limit(params.Limit)
	}

	res := query.Find(&dbEvents)
	if err := res.Error; err != nil {
		return nil, errors.Wrap(err, "failed to get all events")
	}

	ret := make([]models.Event, 0, len(dbEvents))
	for _, dbEvent := range dbEvents {
		event, err := r.GetEventById(ctx, dbEvent.Uid)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get event by id %s", dbEvent.Uid)
		}

		ret = append(ret, event)
	}

	return ret, nil
}

func (r eventRepository) GetSharings(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	var dbEventSharings []db_models.EventSharing

	query := r.db.Model(&db_models.EventSharing{})

	if params.UserID != 0 {
		query = query.Joins("JOIN users on event_sharings.user_id = users.id").
			Joins("JOIN events on event_sharings.event_id = events.id").
			Where("users.uid = ?", params.UserID)
	}

	if params.IsActive.IsDefinedTrue() {
		query = query.Where("events.starts_at >= ?", time.Now().Unix())
	} else if params.IsActive.IsDefinedFalse() {
		query = query.Where("events.starts_at < ?", time.Now().Unix())
	}

	query = query.Where("event_sharings.is_deleted = ?", false)

	query = query.Offset(params.Page * params.Limit)

	if params.Limit != 0 {
		query = query.Limit(params.Limit)
	}

	res := query.Find(&dbEventSharings)
	if err := res.Error; err != nil {
		return nil, errors.Wrap(err, "failed to get user event sharings")
	}

	if len(dbEventSharings) == 0 {
		return nil, nil
	}

	eventIDs := make([]int, 0, len(dbEventSharings))
	for _, sharing := range dbEventSharings {
		eventIDs = append(eventIDs, sharing.EventID)
	}

	var dbEvents []db_models.Event
	res = r.db.Find(&dbEvents, eventIDs)
	if err := res.Error; err != nil {
		return nil, errors.Wrap(err, "failed to get events")
	}

	ret := make([]models.Event, 0, len(dbEvents))
	for _, dbEvent := range dbEvents {
		event, err := r.GetEventById(ctx, dbEvent.Uid)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get event by uid %s", dbEvent.Uid)
		}

		ret = append(ret, event)
	}

	return ret, nil
}

func (r eventRepository) getCategory(name string) (db_models.Category, error) {
	var dbCategory db_models.Category
	res := r.db.Model(&db_models.Category{}).
		Where("name = ?", name).
		Take(&dbCategory)
	if err := res.Error; err != nil {
		return db_models.Category{}, errors.Wrapf(err, "failed to get category by name %s", name)
	}

	return dbCategory, nil
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
			subscriptionEvents, err := r.GetSharings(ctx, models.GetEventParams{
				UserID:   int64(dbSubscription.UserID),
				IsActive: models.DefinedBool(true),
			})
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

func (r eventRepository) GetGroupEvent(ctx context.Context, params models.GetEventParams) ([]models.Event, error) {
	var ret []models.Event
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbGroup db_models.Group
		res := r.db.Take(&dbGroup).Where("group_id = ?", params.GroupId)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get groups events sharings")
		}

		var dbEvents []db_models.Event
		query := r.db.Model(&db_models.Event{}).
			Joins("JOIN groups_events_sharing on groups_events_sharing.event_id = events.id").
			Where("groups_events_sharing.group_id = ?", dbGroup.ID).
			Where("events.is_deleted = ?", false).
			Where("groups_events_sharing.is_deleted = ?", false)

		if params.IsAdmin.IsDefinedTrue() {
			query = query.Where("groups_events_sharing.is_admin = ?", params.IsAdmin.Value)
		}

		if params.IsNeedApprove.IsDefinedTrue() {
			userID := contextlib.GetUser(ctx)

			if userID != int64(dbGroup.UserId) {
				r.logger.Info("try get need approve user no admin %d, %d", userID, dbGroup.UserId)
				return errors.New("try get need approve user no admin")
			}
			query = query.Where("groups_events_sharing.is_need_approve = ?", params.IsNeedApprove.Value)
		}

		if params.Category != "" {
			dbCategory, err := r.getCategory(string(params.Category))
			if err != nil {
				return errors.Wrap(err, "failed to get category")
			}

			query = query.Where("events.category_id = ?", dbCategory.ID)
		}

		if params.City != "" {
			query = query.Where("events.geo LIKE ?", "%"+params.City+"%")
		}

		query = query.Offset(params.Page * params.Limit)

		if params.Limit != 0 {
			query = query.Limit(params.Limit)
		}

		res = query.Find(&dbEvents)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get events group")
		}

		ret = make([]models.Event, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			event, err := r.GetEventById(ctx, dbEvent.Uid)
			if err != nil {
				return errors.Wrapf(err, "failed to get event by id %s", dbEvent.Uid)
			}
			var eventSharing db_models.GroupsEventsSharing

			q := r.db.Model(&eventSharing).Where("group_id = ?", dbGroup.ID).Find(&eventSharing)
			if err := q.Error; err != nil {
				return errors.Wrap(err, "failed to get eventSharing")
			}

			event.GroupInfo.GroupId = params.GroupId
			event.GroupInfo.IsAdmin = eventSharing.IsAdmin

			if !params.IsActive.IsDefined() {
				ret = append(ret, event)
			} else if params.IsActive.IsDefinedTrue() && event.IsActive {
				ret = append(ret, event)
			} else if params.IsActive.IsDefinedFalse() && !event.IsActive {
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
