package postgres

import (
	"context"
	event_repo "github.com/BUSH1997/FrienderAPI/internal/pkg/event/repository/postgres"
	"github.com/pkg/errors"
	"strings"
)

func (r *ImageRepositoryBD) UploadImage(ctx context.Context, uid string, fileName string) error {
	var dbEvent event_repo.Event

	res := r.db.Take(&dbEvent, "uid = ?", uid)
	if err := res.Error; err != nil {
		return errors.Wrap(err, "failed to get event by uid")
	}

	images := strings.Split(dbEvent.Images, ",")
	images = append(images, fileName)
	dbImages := strings.Join(images, ",")

	res = r.db.Model(&event_repo.Event{}).Where("uid = ?", dbEvent.Uid).Update("images", dbImages)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to update images in event, uid %s", uid)
	}

	return nil
}
