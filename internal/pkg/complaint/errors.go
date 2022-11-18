package complaint

import "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"

var (
	ErrAlreadyExists = errors.Typed("complaint_already_exists", "complaint already exists")
)
