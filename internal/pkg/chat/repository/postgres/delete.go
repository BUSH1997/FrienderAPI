package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/chat"
	context2 "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
)

func (r chatRepository) DeleteMessage(ctx context.Context, messageID string) error {
	dbMessage := db_models.Message{}
	res := r.db.Take(&dbMessage, "message_uid = ?", messageID)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to get message %s", messageID)
	}

	eventOwner := db_models.User{}
	res = r.db.Model(&db_models.User{}).
		Joins("JOIN events on events.owner_id = users.id").
		Where("events.uid = ?", dbMessage.EventUID).
		Take(&eventOwner)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to get events owner by event uid %s", dbMessage.EventUID)
	}

	user := context2.GetUser(ctx)
	if user != dbMessage.UserUID && user != int64(eventOwner.Uid) {
		return chat.ErrNotAllowedToDelete
	}

	res = r.db.Model(&db_models.Message{}).
		Where("message_uid = ?", messageID).
		Update("is_deleted", true)
	if err := res.Error; err != nil {
		return errors.Wrapf(err, "failed to delete message %s", messageID)
	}

	return nil
}
