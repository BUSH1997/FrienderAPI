package postgres

import (
	"context"
	context2 "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (r chatRepository) GetMessages(ctx context.Context, opts models.GetMessageOpts) ([]models.Message, error) {
	var messages []models.Message

	offset := opts.Page * opts.Limit

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbMessages []db_models.Message

		res := r.db.Model(&db_models.Message{}).
			Joins("JOIN events on messages.event_id = events.id").
			Order("time_created DESC").
			Offset(offset).
			Limit(opts.Limit).
			Find(&dbMessages, "events.uid = ?", opts.EventID)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get messages")
		}

		messages = make([]models.Message, 0, len(dbMessages))
		for _, dbMessage := range dbMessages {
			message := models.Message{
				UserID:      dbMessage.UserUID,
				Text:        dbMessage.Text,
				EventID:     dbMessage.EventUID,
				TimeCreated: dbMessage.TimeCreated,
			}

			messages = append(messages, message)
		}

		//for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		//	messages[i], messages[j] = messages[j], messages[i]
		//}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return messages, nil
}

func (r chatRepository) GetChats(ctx context.Context) ([]models.Chat, error) {
	var chats []models.Chat

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvents []db_models.Event
		res := r.db.Model(&db_models.Event{}).
			Joins("JOIN event_sharings on event_sharings.event_id = events.id").
			Joins("JOIN users on event_sharings.user_id = users.id").
			Find(&dbEvents, "users.uid = ? AND event_sharings.is_deleted = ?", context2.GetUser(ctx), false)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user events")
		}

		chats = make([]models.Chat, 0, len(dbEvents))
		for _, dbEvent := range dbEvents {
			chat := models.Chat{
				EventUID:    dbEvent.Uid,
				EventTitle:  dbEvent.Title,
				EventAvatar: dbEvent.AvatarUrl,
			}

			chats = append(chats, chat)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return chats, nil
}
