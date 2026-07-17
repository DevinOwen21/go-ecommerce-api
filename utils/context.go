package utils

import "context"

type contextKey string

type UserContext struct {
	UserID string
	Role   string
}

const userContextKey contextKey = "userContext"

func SetUserContext(ctx context.Context, user UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserContext(ctx context.Context) (UserContext, bool) {
	user, ok := ctx.Value(userContextKey).(UserContext)
	return user, ok
}
