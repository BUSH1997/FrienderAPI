package context

import (
	"context"
	"fmt"
	"math/rand"
)

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

type requestIDKey struct{}

func GetRequestID(ctx context.Context) string {
	reqID := ctx.Value(requestIDKey{})
	if reqID == nil {
		return ""
	}

	return reqID.(string)
}

func GetOrGenerateRequestID(ctx context.Context) string {
	reqID := ctx.Value(requestIDKey{})
	if reqID == nil {
		return GenerateRequestID()
	}

	return reqID.(string)
}

func SetRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, reqID)
}

func GenerateRequestID() string {
	return fmt.Sprintf("%016x", rand.Int())
}
