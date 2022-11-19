package convert

import (
	api_errors "github.com/BUSH1997/FrienderAPI/internal/api/errors"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/complaint"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/event"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/profile"
	"github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"
)

func DeliveryError(err error) api_errors.ApiError {
	if err == nil {
		return nil
	}

	if errors.Is(err, event.ErrEmptyUpdate) {
		return api_errors.ValidationError{
			Cause:                 err,
			Explain:               "empty update not allowed",
			ValidationErrorFields: api_errors.ValidationErrorFields{},
		}
	}

	if errors.Is(err, event.ErrEventIsSpam) {
		return api_errors.ForbiddenError{
			Cause:  err,
			Reason: api_errors.EventIsSpam,
		}
	}

	if errors.Is(err, event.ErrAlreadyExists) {
		return api_errors.AlreadyExistsError{
			Cause:  err,
			Entity: "event",
		}
	}

	if errors.Is(err, event.ErrNotExists) {
		return api_errors.NotFoundError{
			Cause: err,
			Key:   "event",
		}
	}

	if errors.Is(err, event.ErrNoFullAccess) {
		return api_errors.ForbiddenError{
			Cause:  err,
			Reason: "no_edit_access_to_event",
		}
	}

	if errors.Is(err, event.ErrNoWriteAccess) {
		return api_errors.ForbiddenError{
			Cause:  err,
			Reason: "no_write_access_to_event",
		}
	}

	if errors.Is(err, profile.ErrNotFound) {
		return api_errors.NotFoundError{
			Cause: err,
			Key:   "user",
		}
	}

	if errors.Is(err, profile.ErrAlreadyExists) {
		return api_errors.AlreadyExistsError{
			Cause:  err,
			Entity: "user",
		}
	}

	if errors.Is(err, profile.ErrEmptyUpdate) {
		return api_errors.ValidationError{
			Cause:   err,
			Explain: "empty update not allowed",
			ValidationErrorFields: api_errors.ValidationErrorFields{
				"user_update_input": api_errors.FieldInvalid,
			},
		}
	}

	if errors.Is(err, complaint.ErrAlreadyExists) {
		return api_errors.AlreadyExistsError{
			Cause:  err,
			Entity: "complaint",
		}
	}

	return api_errors.InternalError{Cause: errors.New("internal server error")}
}
