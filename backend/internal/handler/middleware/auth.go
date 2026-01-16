package middleware

import "context"

type ctxKeyUserID struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxKeyUserID{}, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxKeyUserID{})
	userID, ok := v.(string)

	return userID, ok
}
