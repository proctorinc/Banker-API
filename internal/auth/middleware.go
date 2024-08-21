package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/db"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// A stand-in for our database backed user object
type User struct {
	Name    string
	IsAdmin bool
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(db db.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")

		// Deny unauthorized users
		if token == "" || err != nil || ctx == nil {
			ctx.Next()
			return
		}

		// Validate user Id from token
		userId, err := validateAndGetUserID(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth token supplied"})
			ctx.Abort()
			return
		}

		// get the user from the database
		user, err := db.GetUser(ctx, userId)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth token supplied"})
			ctx.Abort()
			return
		}

		// Add user to request context
		ctx.Set(UserContextKey, user)
		ctx.Next()
	}
}
