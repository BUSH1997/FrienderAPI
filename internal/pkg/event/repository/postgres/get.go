package postgres

import (
	"context"
	"fmt"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func (r eventRepository) GetAllPublic(ctx context.Context) ([]models.Event, error) {
	var dbEvents []Event

	res := r.db.Find(&dbEvents, "is_public = ?", true)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err := res.Error; err != nil {
		return nil, errors.Wrap(err, "failed to check user")
	}

	events := make([]models.Event, 0, len(dbEvents))
	for _, dbEvent := range dbEvents {
		event := models.Event{
			Uid:      dbEvent.Uid,
			Title:    dbEvent.Title,
			StartsAt: dbEvent.StartsAt,
			IsPublic: dbEvent.IsPublic,
		}

		events = append(events, event)
	}

	return events, nil
}

func (r eventRepository) GetEventById(ctx context.Context, id string) (models.Event, error) {
	var dbEvent Event
	res := r.db.Take(&dbEvent, "uid = ?", id)
	if err := res.Error; err != nil {
		return models.Event{}, errors.Wrap(err, "failed to get event by id")
	}

	var dbUser User
	res = r.db.Take(&dbUser, "id = ?", dbEvent.Owner)
	if err := res.Error; err != nil {
		return models.Event{}, errors.Wrap(err, "failed to get owner id")
	}

	var dbCategory Category
	res = r.db.Take(&dbCategory, "id = ?", dbEvent.Category)
	if err := res.Error; err != nil {
		return models.Event{}, errors.Wrap(err, "failed to get category id")
	}

	var dbEventSharings []EventSharing

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

	var dbMembers []User
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
		TimeCreated: dbEvent.TimeCreated,
		TimeUpdated: dbEvent.TimeUpdated,
		Author:      dbUser.Uid,
		StartsAt:    dbEvent.StartsAt,
		IsGroup:     dbEvent.IsGroup,
		IsPublic:    dbEvent.IsPublic,
		Category:    models.Category(dbCategory.Name),
	}

	event.Members = members

	longitude, err := strconv.ParseFloat(strings.Split(dbEvent.Geo, ",")[0], 32)
	if err != nil {
		fmt.Println(err, "POP")
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
	var dbEvents []Event
	res := r.db.Find(&dbEvents)
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

func (r eventRepository) GetUserEvents(ctx context.Context, id string) ([]models.Event, error) {
	return nil, nil
}
