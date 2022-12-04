package chat

import "github.com/BUSH1997/FrienderAPI/internal/pkg/tools/errors"

var (
	ErrNotAllowedToDelete = errors.Typed("cannot_delete_message", "user cannot delete message")
	UnexpectedMessageType = errors.Typed("unexpected_message_type", "unexpected message type")
)
