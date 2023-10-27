package internalcontext

import (
	"context"

	"github.com/google/uuid"
)

var userIDKey = struct{}{}

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) uuid.UUID {
	return ctx.Value(userIDKey).(uuid.UUID)
}
