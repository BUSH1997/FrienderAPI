package profile

import "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"

var (
	ErrNotFound      = errors.Typed("user_not_found", "user not found")
	ErrEmptyUpdate   = errors.Typed("user_empty_update", "user empty update")
	ErrAlreadyExists = errors.Typed("user_already_exists", "user already exists")
)
