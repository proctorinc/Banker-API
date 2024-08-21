package auth

import (
	"context"

	"github.com/google/uuid"
)

const userContextKey = "user"

func GetCurrentUser(ctx context.Context) *User {
	raw, _ := ctx.Value(userContextKey).(*User)
	return raw
}

func Login(email string, password string) (bool, error) {
	return true, nil
}

func isAuthenticated(ctx context.Context) bool {
	return GetCurrentUser(ctx) != nil
}

func validateAndGetUserID(token string) (uuid.UUID, error) {
	return uuid.Parse(token)
}
