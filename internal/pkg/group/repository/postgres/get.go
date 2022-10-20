package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
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
