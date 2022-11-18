package postgres

import (
	"context"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	db_models "github.com/BUSH1997/FrienderAPI/internal/pkg/postgres/models"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
	"gorm.io/gorm"
	"strings"
)

func (r profileRepository) GetOneProfile(ctx context.Context, userID int64) (models.Profile, error) {
	ctx = r.logger.WithCaller(ctx)

	currentStatus, err := r.getUserCurrentStatus(ctx, userID)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get profile current status")
	}

	awards, err := r.getUserAwards(ctx, userID)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get user awards")
	}

	profile := models.Profile{
		ProfileStatus: currentStatus,
		Awards:        awards,
		CanBeReported: true,
	}

	initiator := contextlib.GetUser(ctx)

	var dbInitiator db_models.User
	res := r.db.Take(&dbInitiator, "uid = ?", initiator)
	if err := res.Error; err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get current user id")
	}

	var dbComplaint db_models.Complaint
	res = r.db.Model(&db_models.Complaint{}).
		Where("item = ?", "user").
		Where("item_uid = ?", userID).
		Where("initiator = ?", dbInitiator.ID).
		Take(&dbComplaint)
	if err := res.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Profile{}, errors.Wrap(err, "failed to get complaint")
	}
	if userID == initiator || !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		profile.CanBeReported = false
	}

	return profile, nil
}

func (r profileRepository) CheckUserExists(ctx context.Context, user int64) (bool, error) {
	ctx = r.logger.WithCaller(ctx)

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

func (r profileRepository) GetAllUserStatuses(ctx context.Context, id int64) ([]models.Status, error) {
	ctx = r.logger.WithCaller(ctx)

	var statuses []models.Status

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var dbStatuses []db_models.Status

		res := r.db.Model(&db_models.Status{}).
			Joins("JOIN unlocked_statuses on unlocked_statuses.status_id = statuses.id").
			Joins("JOIN users on users.id = unlocked_statuses.user_id").
			Find(&dbStatuses, "users.uid = ?", id)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		if err := res.Error; err != nil {
			return errors.Wrap(err, "failed to get user statuses")
		}

		statuses = make([]models.Status, 0, len(dbStatuses))
		for _, dbStatus := range dbStatuses {
			status := models.Status{
				Title:    dbStatus.Title,
				Id:       dbStatus.UID,
				IsLocked: false,
			}

			statuses = append(statuses, status)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make transaction")
	}

	return statuses, nil
}

func (r profileRepository) getUserCurrentStatus(ctx context.Context, id int64) (models.Status, error) {
	ctx = r.logger.WithCaller(ctx)

	var status models.Status

	var dbUser db_models.User
	res := r.db.Take(&dbUser, "uid = ?", id)
	if err := res.Error; err != nil {
		return models.Status{}, errors.Wrap(err, "failed to get user")
	}

	var dbStatus db_models.Status
	res = r.db.Take(&dbStatus, "user_id = ?", dbUser.ID)
	if err := res.Error; err != nil {
		return models.Status{}, errors.Wrapf(err, "failed to get status")
	}

	status.Title = dbStatus.Title
	status.Id = dbStatus.UID
	status.IsLocked = false

	return status, nil
}

func (r profileRepository) getUserAwards(ctx context.Context, id int64) ([]models.Award, error) {
	var awards []models.Award
	var dbAwards []db_models.Award

	res := r.db.Model(&db_models.Award{}).
		Joins("JOIN unlocked_awards on unlocked_awards.award_id = awards.id").
		Joins("JOIN users on users.id = unlocked_awards.user_id").
		Where("users.uid = ?", id).
		Find(&dbAwards)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err := res.Error; err != nil {
		return nil, errors.Wrap(err, "failed to get user awards")
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

	return awards, nil
}
