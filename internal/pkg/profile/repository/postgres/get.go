package postgres

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"
)

func (r profileRepository) CheckUserExists(ctx context.Context, user int64) (bool, error) {
	userExists := true
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbUser db_models.User
		res := r.db.Take(&dbUser, "uid = ?", user)
		if err := res.Error; errors.Is(err, gorm.ErrRecordNotFound) {
			userExists = false
			return nil
		}
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user by id")
		}

		return nil
	})
	if err != nil {
		return false, errors.Wrap(err, "failed to make transaction")
	}

	return userExists, nil
}

func (r profileRepository) GetSubscribe(cxt context.Context, userId int64) ([]models.SubscribeType, error) {
	var result []models.SubscribeType
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var currentUser db_models.User
		res := r.db.Take(&currentUser, "uid = ?", userId)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get subscribe")
		}

		var dbSubscribers []db_models.SubscribeProfileSharing

		res = r.db.Find(&dbSubscribers, "user_id = ?", currentUser.ID)
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get subscribe")
		}

		for _, dbSubscribe := range dbSubscribers {
			var user db_models.User
			res := r.db.Take(&user, "id = ?", dbSubscribe.ProfileId)
			if err := res.Error; err != nil {
				return errors.Wrap(err, "failed to get subscribe")
			}

			subscribeType := models.SubscribeType{
				Id:      int64(user.Uid),
				IsGroup: user.IsGroup,
			}

			result = append(result, subscribeType)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return result, nil

}

func (r profileRepository) GetCities(ctx context.Context) ([]string, error) {
	cities := make([]string, 0)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbEvents []db_models.Event
		res := r.db.Find(&dbEvents, "geo LIKE ?", "%;;%;;_%")
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get all categories")
		}
		mapCities := make(map[string]bool)
		for _, v := range dbEvents {
			geoArray := strings.Split(v.Geo, ";;")
			city := strings.Split(geoArray[2], ",")
			mapCities[city[0]] = true
		}

		for k, _ := range mapCities {
			cities = append(cities, k)
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction GetAllCategories")
	}

	return cities, nil
}
