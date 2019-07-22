package userclient

import "context"

type ctxKey int

const (
	userCtxKey ctxKey = iota
	userTokenCtxKey
)

// GetCurrentUserFromContext return user from context or nil
func GetCurrentUserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value(userCtxKey).(*User); ok {
		return user
	}

	return nil
}

// ContextWithUser add user to context
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

// GetTokenFromContext return token from context, empty string if token not exists
func GetTokenFromContext(ctx context.Context) string {
	if token, ok := ctx.Value(userTokenCtxKey).(string); ok {
		return token
	}

	return ""
}

// ContextWithToken add user token to context
func ContextWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, userTokenCtxKey, token)
}
