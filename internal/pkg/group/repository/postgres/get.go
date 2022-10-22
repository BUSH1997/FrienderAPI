package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func (gr *groupRepository) GetAdministeredGroupByUserId(ctx context.Context, userId int) ([]models.Group, error) {
	var ret []models.Group

	err := gr.db.Transaction(func(tx *gorm.DB) error {
		var dbGroups []models.Group

		res := gr.db.Find(&dbGroups).Where("user_id = ?", userId)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get all events")
		}

		for _, dbGroup := range dbGroups {
			group := models.Group{
				GroupId: dbGroup.GroupId,
				UserId:  dbGroup.UserId,
			}

			ret = append(ret, group)
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "[GetAdministeredGroupByUserId] failed to make transaction")
	}

	return ret, nil
}

func (gr *groupRepository) CheckIfAdmin(ctx context.Context, userId int, groupId int) (bool, error) {
	isAdmin := true
	err := gr.db.Transaction(func(tx *gorm.DB) error {
		var dbGroup db_models.Group
		res := gr.db.Take(&dbGroup, "group_id = ? and user_id = ?", groupId, userId)
		if err := res.Error; errors.Is(err, gorm.ErrRecordNotFound) {
			isAdmin = false
			return nil
		}
		if err := res.Error; err != nil {
			return errors.Wrap(err, "[CheckIfAdmin] failed to get group by user_id and group_id")
		}

		return nil
	})
	if err != nil {
		return false, errors.Wrap(err, "[CheckIfAdmin] failed to make transaction")
	}

	return isAdmin, nil
}
