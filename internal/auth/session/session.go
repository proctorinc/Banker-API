package session

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/db"
)

type Session struct {
	Writer     gin.ResponseWriter
	IsLoggedIn bool
	User       *db.User
}

const sessionContextKey = "user-session"

func GetSession(ctx context.Context) *Session {
	return ctx.Value(sessionContextKey).(*Session)
}

func SetSession(ctx *gin.Context, session *Session) {
	newCtx := context.WithValue(ctx.Request.Context(), sessionContextKey, session)
	ctx.Request = ctx.Request.WithContext(newCtx)
}
