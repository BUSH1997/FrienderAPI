package event

import "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"

var (
	ErrNotExists         = errors.Typed("event_not_exists", "event doesn't exist")
	ErrNotFoundInIndex   = errors.Typed("event_not_found_in_index", "event is not found in revindex")
	ErrAlreadyExists     = errors.Typed("event_already_exists", "event already exists")
	ErrEventIsSpam       = errors.Typed("event_is_spam", "event is spam")
	ErrNoFullAccess      = errors.Typed("event_no_full_access", "no full access to event")
	ErrNoWriteAccess     = errors.Typed("event_no_write_access", "no write access to event")
	ErrEmptyUpdate       = errors.Typed("event_empty", "empty update")
	ErrNoDeleteAccess    = errors.Typed("event_no_delete_access", "no delete access to event")
	ErrNoAccessForBanned = errors.Typed("event_no_access_for_banned", "no access to event for banned user")
)
