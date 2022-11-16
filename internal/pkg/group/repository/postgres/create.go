package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/gorm"
)

func (gr *groupRepository) Create(ctx context.Context, group models.GroupInput) error {
	ctx = gr.logger.WithCaller(ctx)

	err := gr.db.Transaction(func(tx *gorm.DB) error {
		dbGroup := db_models.Group{
			UserId:          group.UserId,
			GroupId:         group.GroupId,
			AllowUserEvents: group.AllowUserEvents,
		}
		res := gr.db.Create(&dbGroup)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "[CreateGroup] failed to create group")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "[CreateGroup] failed to make transaction ")
	}

	return nil
}
