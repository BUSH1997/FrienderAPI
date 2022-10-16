package context

import "context"

type userKey struct{}

func SetUser(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userKey{}, userID)
}

func GetUser(ctx context.Context) int64 {
	userID := ctx.Value(userKey{})
	if userID == nil {
		return 0
	}

	return userID.(int64)
}
