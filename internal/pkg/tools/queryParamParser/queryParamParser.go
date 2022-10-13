package queryParamParser

import (
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"net/url"
	"strconv"
)

const (
	DefaultPageLimit = 50
	DefaultPageStart = 0
)

func ParseGetAllEvents(values url.Values) (event.FilterGetAll, error) {
	var result event.FilterGetAll
	result.IsSubscribe = false
	result.Page = DefaultPageStart
	result.Limit = DefaultPageLimit

	isSubString := values.Get("is_sub")
	if isSubString != "" {
		isSub, err := strconv.ParseBool(isSubString)
		if err != nil {
			return event.FilterGetAll{}, err
		}
		result.IsSubscribe = isSub
	}

	isActiveString := values.Get("is_active")
	if isActiveString != "" {
		isActive, err := strconv.ParseBool(isActiveString)
		if err != nil {
			return event.FilterGetAll{}, err
		}
		result.IsActive = isActive
	}

	isOwnerString := values.Get("is_owner")
	if isOwnerString != "" {
		isOwner, err := strconv.ParseBool(isOwnerString)
		if err != nil {
			return event.FilterGetAll{}, err
		}
		result.IsOwner = isOwner
	}

	userString := values.Get("user")
	if userString != "" {
		userId, err := strconv.Atoi(userString)
		if err != nil {
			return event.FilterGetAll{}, err
		}
		result.User = userId
	}

	pageLimitString := values.Get("limit")
	if pageLimitString != "" {
		pageLimit, err := strconv.Atoi(pageLimitString)
		if err != nil {
			return event.FilterGetAll{}, err
		}
		result.Limit = pageLimit
	}

	pageString := values.Get("page")
	if pageString != "" {
		page, err := strconv.Atoi(pageString)
		if err != nil {
			return event.FilterGetAll{}, err
		}
		result.Page = page
	}

	return result, nil
}
