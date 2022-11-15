package usecase

import (
	"context"
	"github.com/BUSH1997/FrienderAPI/cmd/public_sync/client/vk"
	contextlib "github.com/BUSH1997/FrienderAPI/internal/pkg/context"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/models"
	"github.com/pkg/errors"
	"strconv"
)

func (uc *UseCase) GetOneProfile(ctx context.Context, userID int64) (models.Profile, error) {
	ctx = uc.Logger.WithCaller(ctx)

	currentStatus, err := uc.statusRepository.GetUserCurrentStatus(ctx, userID)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get profile status")
	}

	activeEvents, err := uc.eventRepository.GetSharings(ctx, models.GetEventParams{
		UserID:   userID,
		IsActive: models.DefinedBool(true),
	})
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get user active events")
	}

	visitedEvents, err := uc.eventRepository.GetSharings(ctx, models.GetEventParams{
		UserID:   userID,
		IsActive: models.DefinedBool(false),
	})
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get user visited events")
	}

	awards, err := uc.awardRepository.GetUserAwards(ctx, userID)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "failed to get user awards")
	}

	profile := models.Profile{
		ProfileStatus: currentStatus,
		Awards:        awards,
		ActiveEvents:  activeEvents,
		VisitedEvents: visitedEvents,
	}

	return profile, nil
}

func (uc *UseCase) GetAllProfileStatuses(ctx context.Context) ([]models.Status, error) {
	ctx = uc.Logger.WithCaller(ctx)

	userID := contextlib.GetUser(ctx)

	statuses, err := uc.statusRepository.GetAllUserStatuses(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all profile statuses")
	}

	return statuses, nil
}

func (uc *UseCase) GetSubscribe(ctx context.Context, userId int64) (models.Subscriptions, error) {
	ctx = uc.Logger.WithCaller(ctx)

	subscribe, err := uc.profileRepository.GetSubscribe(ctx, userId)
	if err != nil {
		uc.Logger.WithError(err).Errorf("[GetSubscribe] failed get subscribe")
		return models.Subscriptions{
			Groups: []int64{},
			Users:  []int64{},
		}, nil
	}

	Users := make([]int64, 0)
	Groups := make([]int64, 0)
	for _, v := range subscribe {
		if v.IsGroup {
			Groups = append(Groups, int64(v.Id))
		} else {
			Users = append(Users, int64(v.Id))
		}
	}

	return models.Subscriptions{
		Groups: Groups,
		Users:  Users,
	}, nil
}

func (uc *UseCase) GetFriends(ctx context.Context, userId int64) ([]int64, error) {
	ctx = uc.Logger.WithCaller(ctx)

	getFriendsFormData := map[string]string{
		"access_token": "vk1.a.3v18zK0yJZRszF9FRAvhVhACDcDYPqZeeEkaehZ0k-qli2EIioZif1R4mI1cfQuwxH7cqLXG2JmDGHcf4AiTma5MpwGnhyZ3FBWjMbLqlbvCjRk1AbK8_7oWxO0DZBRySBUh2XDWCtXY6SVRRl4gDq07_U3IC-IdASY5nzcVTgZ7-qoib3C8fhoU-6I1U7-e",
		"user_id":      strconv.Itoa(int(userId)),
		"v":            "5.131",
	}
	respFriends := vk.GetFriendsResponse{
		DownloadLimitBytes: 2000000000,
	}

	err := uc.httpClient.PerformRequest(ctx, vk.GetRequestWithBody{
		GetRequest: vk.GetRequest{
			RequestURL: "https://api.vk.com/method/friends.get",
		},
		FormData: getFriendsFormData,
	}, &respFriends)
	if err != nil {
		uc.Logger.WithCtx(ctx).WithError(err).Errorf("failed to perform request")
		return []int64{}, errors.Wrap(err, "failed to perform request")
	}

	result := make([]int64, 0)
	for _, v := range respFriends.VkFriendsData.Ids {
		isExist, err := uc.profileRepository.CheckUserExists(ctx, v)
		if err != nil {
			uc.Logger.WithCtx(ctx).WithError(err).Errorf("failed to check if user exists")
			return []int64{}, errors.Wrap(err, "failed to check if user exists")
		}
		if isExist {
			result = append(result, v)
		}
	}

	return result, nil
}

func (uc *UseCase) GetCities(ctx context.Context) ([]string, error) {
	ctx = uc.Logger.WithCaller(ctx)

	cities, err := uc.profileRepository.GetCities(ctx)
	if err != nil {
		uc.Logger.WithCtx(ctx).WithError(err).Errorf("failed to get cities in usecase")
		return cities, err
	}

	return cities, nil
}
