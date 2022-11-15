package postgres

import (
	"context"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

func (r *eventRepository) UploadImage(ctx context.Context, uid string, link string) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvent db_models.Event

		res := r.db.Take(&dbEvent, "uid = ?", uid)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get event by uid")
		}

		res = r.db.Model(&db_models.Event{}).
			Where("uid = ?", dbEvent.Uid).Update("images", buildImageLink(dbEvent.Images, link))
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update images in event, uid %s", uid)
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}

func buildImageLink(dbImages string, link string) string {
	if dbImages == "" {
		return link
	}

	images := strings.Split(dbImages, ",")
	images = append(images, link)

	return strings.Join(images, ",")
}

func (r *eventRepository) UploadAvatar(ctx context.Context, uid string, link string, vkId string) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvent db_models.Event

		res := r.db.Take(&dbEvent, "uid = ?", uid)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get event by uid")
		}

		res = r.db.Model(&db_models.Event{}).
			Where("uid = ?", dbEvent.Uid).Updates(map[string]interface{}{
			"avatar_url":   link,
			"avatar_vk_id": vkId,
		})
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update images in event, uid %s", uid)
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}
