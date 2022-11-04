package postgres

import (
	"context"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r eventRepository) Delete(ctx context.Context, event string, groupInfo models.GroupInfo) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvent db_models.Event
		res := r.db.Take(&dbEvent, "uid = ?", event)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get event by uid %s", event)
		}

		userID := contextlib.GetUser(ctx)

		var dbUser db_models.User
		res = r.db.Take(&dbUser, "uid = ?", userID)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get user by uid %d", userID)
		}

		res = r.db.Model(&db_models.User{}).Where("id = ?", dbUser.ID).
			Update("created_events", dbUser.CreatedEventsCount-1)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to update user created events")
		}

		if groupInfo.GroupId == 0 {
			if dbEvent.Owner != int(dbUser.ID) {
				return errors.New("user is not events owner, cannot delete")
			}
		} else {
			//надо проверить является ли юзер админом
			var groups []db_models.Group
			res = r.db.Find(&groups, "user_id = ?", userID)
			if err := res.Error; err != nil {
				return errors.Wrap(err, "failed to get admin group")
			}

			checkAdmin := false
			for _, group := range groups {
				if int64(group.GroupId) == groupInfo.GroupId {
					checkAdmin = true
				}
			}

			if !checkAdmin {
				return errors.New("user is not admin, cannot delete")
			}
		}

		dbEventSharing := db_models.EventSharing{}
		res = r.db.Take(&dbEventSharing, "event_id = ? AND user_id = ?", dbEvent.ID, dbUser.ID)

		currentPriority := dbEventSharing.Priority

		res = r.db.Model(&db_models.EventSharing{}).
			Where("user_id = ? AND priority > ?", dbUser.ID, currentPriority).
			Update("priority", gorm.Expr("event_sharings.priority + ?", -1))
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharings priority")
		}

		res = r.db.Model(&db_models.EventSharing{}).
			Where("event_id = ?", dbEvent.ID).
			Update("is_deleted", true)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event sharings")
		}

		res = r.db.Model(&db_models.Event{}).
			Where("id = ?", dbEvent.ID).
			Update("is_deleted", true)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update event")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func (r eventRepository) DeletePhotos(ctx context.Context, links []string, uid string) error {
	var dbUser db_models.User
	res := r.db.Take(&dbUser, "uid = ?", contextlib.GetUser(ctx))
	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to get user by id")
	}

	var dbEvent db_models.Event
	res = r.db.Take(&dbEvent, "uid = ?", uid)
	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to get event by uid")
	}

	var dbEventSharing db_models.EventSharing
	res = r.db.Model(&db_models.EventSharing{}).
		Where("event_id = ?", dbEvent.ID).
		Where("user_id = ?", dbUser.ID).
		Take(&dbEventSharing)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to get event sharing")
	}

	existPhotos := []string(dbEventSharing.Photos)
	existPhotoMap := make(map[string]bool)
	for _, existPhoto := range existPhotos {
		existPhotoMap[existPhoto] = true
	}

	for _, link := range links {
		delete(existPhotoMap, link)
	}

	leftPhotos := make([]string, 0, len(existPhotoMap))
	for existPhoto := range existPhotoMap {
		leftPhotos = append(leftPhotos, existPhoto)
	}

	res = r.db.Model(&db_models.EventSharing{}).
		Where("event_id = ?", dbEvent.ID).
		Where("user_id = ?", dbUser.ID).
		Update("photos", leftPhotos)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to upload photos in event, uid %s", uid)
	}

	return nil
}
