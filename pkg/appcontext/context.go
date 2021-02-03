package appcontext

import (
	"context"
)

type contextKey int

const (
	authSubjectKey contextKey = iota
)

// WithAuthSubject adds the Authorization JWT subject to the context.
func WithAuthSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, authSubjectKey, subject)
}

// GetAuthSubject gets the Authorization JWT subject from the context.
func GetAuthSubject(ctx context.Context) (string, bool) {
	tokenStr, ok := ctx.Value(authSubjectKey).(string)
	return tokenStr, ok
}
