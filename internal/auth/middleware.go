package auth

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/auth/session"
	"github.com/proctorinc/banker/internal/auth/token"
	"github.com/proctorinc/banker/internal/db"
)

type contextKey struct {
	name string
}

func Middleware(db db.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqSession := session.Session{
			Writer:     ctx.Writer,
			IsLoggedIn: false,
			User:       nil,
		}

		session.SetSession(ctx, &reqSession)

		authToken, err := token.GetAuthToken(ctx)

		// Deny unauthorized users
		if authToken.IsEmpty() || err != nil || ctx == nil {
			log.Println("No auth token supplied")
			ctx.Next()
			return
		}

		log.Printf("AuthToken: %s", authToken.Value)

		userId, err := authToken.GetUserId()

		if err != nil {
			log.Println("Invalid token provided #1")
			ctx.Next()
			return
		}

		user, err := db.GetUser(ctx, userId)

		if err != nil {
			log.Println("Invalid token provided #2")
			ctx.Next()
			return
		}

		// Add user to reques context
		reqSession.IsLoggedIn = true
		reqSession.User = &user

		ctx.Next()
	}
}
