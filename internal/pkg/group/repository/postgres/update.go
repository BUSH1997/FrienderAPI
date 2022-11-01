package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (gr *groupRepository) Update(ctx context.Context, group models.GroupInput) error {
	err := gr.db.Transaction(func(tx *gorm.DB) error {
		res := gr.db.Model(&db_models.Group{}).
			Where("group_id = ?", group.GroupId).
			Updates(map[string]interface{}{
				"user_id":           group.UserId,
				"group_id":          group.GroupId,
				"allow_user_events": group.AllowUserEvents,
			})
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update group, groupID = %d", group.GroupId)
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}
