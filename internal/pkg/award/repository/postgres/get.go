package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/gorm"
)

func (r awardRepository) GetUserAwards(ctx context.Context, id int64) ([]models.Award, error) {
	var awards []models.Award
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbAwards []db_models.Award

		res := r.db.Model(&db_models.Award{}).
			Joins("JOIN unlocked_awards on unlocked_awards.award_id = awards.id").
			Joins("JOIN users on users.id = unlocked_awards.user_id").
			Find(&dbAwards, "users.uid = ?", id)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user awards")
		}

		awards = make([]models.Award, 0, len(dbAwards))
		for _, dbAward := range dbAwards {
			award := models.Award{
				Image:       dbAward.Image,
				Name:        dbAward.Name,
				Description: dbAward.Description,
				IsLocked:    false,
			}

			awards = append(awards, award)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return awards, nil
}
