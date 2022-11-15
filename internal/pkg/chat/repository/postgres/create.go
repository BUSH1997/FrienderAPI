package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r chatRepository) CreateMessage(ctx context.Context, message models.Message) error {
	ctx = r.logger.WithCaller(ctx)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbUser db_models.User
		res := r.db.Take(&dbUser, "uid = ?", message.UserID) //message.UserID)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user")
		}

		var dbEvent db_models.Event
		res = r.db.Take(&dbEvent, "uid = ?", message.EventID)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get event")
		}

		dbMessage := db_models.Message{
			UserID:      int(dbUser.ID),
			UserUID:     int64(dbUser.Uid),
			TimeCreated: message.TimeCreated,
			Text:        message.Text,
			EventID:     int(dbEvent.ID),
			EventUID:    dbEvent.Uid,
		}

		res = r.db.Create(&dbMessage)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to create message")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}
