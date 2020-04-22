package oauth2

import (
	"context"
)

var userKey struct{}

type User struct {
	Issuer  string `json:"iss"`
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

func AuthenticatedUser(ctx context.Context) *User {
	if user, ok := ctx.Value(&userKey).(*User); ok {
		return user
	}
	return nil
}

func WithAuthenticatedUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, &userKey, user)
}
