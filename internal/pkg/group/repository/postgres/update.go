package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/gorm"
)

func (gr *groupRepository) Update(ctx context.Context, group models.GroupInput) error {
	ctx = gr.logger.WithCaller(ctx)

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

func (gr *groupRepository) ApproveEvent(ctx context.Context, eventApproveInfo models.ApproveEvent) error {
	ctx = gr.logger.WithCaller(ctx)

	err := gr.db.Transaction(func(tx *gorm.DB) error {
		var dbEvent db_models.Event
		res := gr.db.Take(&dbEvent, "uid = ?", eventApproveInfo.EventUid)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get event")
		}

		var dbGroup db_models.Group
		res = gr.db.Take(&dbGroup, "group_id = ?", eventApproveInfo.GroupId)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to get group")
		}

		var approveUpdate map[string]interface{}

		if eventApproveInfo.Approve {
			approveUpdate = map[string]interface{}{
				"is_need_approve": false,
				"is_admin":        true,
			}
		} else {
			approveUpdate = map[string]interface{}{
				"is_deleted": true,
			}
		}

		res = gr.db.Model(&db_models.GroupsEventsSharing{}).
			Where("group_id = ?", dbGroup.ID).
			Where("event_id = ?", dbEvent.ID).
			Updates(approveUpdate)
		if err := res.Error; err != nil {
			return errors.Wrapf(err, "failed to update GroupsEventsSharing")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to make transaction")
	}

	return nil
}
