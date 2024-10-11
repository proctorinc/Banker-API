package auth

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/auth/jwt"
	"github.com/proctorinc/banker/internal/auth/session"
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

		tokenString, err := session.GetAuthTokenString(ctx)

		if err != nil {
			log.Printf("no jwt token present: %v", err)
			ctx.Next()
			return
		}

		jwtToken, err := jwt.VerifyJwt(tokenString)

		if err != nil {
			log.Printf("error validating JWT token: %v", err)
			ctx.Next()
			return
		}

		if time.Now().After(jwtToken.ExpiresAt) {
			log.Println("jwt token has expired")
			ctx.Next()
			return
		}

		user, err := db.GetUser(ctx, jwtToken.UserId)

		if err != nil {
			ctx.Next()
			return
		}

		// Add user to request context
		reqSession.IsLoggedIn = true
		reqSession.User = &user

		ctx.Next()
	}
}
