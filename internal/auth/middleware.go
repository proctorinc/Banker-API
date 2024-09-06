package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/auth/token"
	"github.com/proctorinc/banker/internal/db"
)

type contextKey struct {
	name string
}

func Middleware(db db.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken, err := token.GetAuthToken(ctx)

		// Deny unauthorized users
		if authToken.IsEmpty() || err != nil || ctx == nil {
			ctx.Next()
			return
		}

		userId, err := authToken.GetUserId()

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth token supplied"})
			ctx.Abort()
			return
		}

		user, err := db.GetUser(ctx, userId)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth token supplied"})
			ctx.Abort()
			return
		}

		// Add user to request context
		SetAuthenticatedUser(ctx, user)
		ctx.Next()
	}
}
