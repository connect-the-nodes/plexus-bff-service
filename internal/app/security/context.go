package security

import "context"

type contextKey string

const authContextKey contextKey = "auth"

type Authentication struct {
	Claims      map[string]any
	Authorities []string
	Token       string
}

func WithAuthentication(ctx context.Context, auth *Authentication) context.Context {
	return context.WithValue(ctx, authContextKey, auth)
}

func FromContext(ctx context.Context) *Authentication {
	value, _ := ctx.Value(authContextKey).(*Authentication)
	return value
}
