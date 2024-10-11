package session

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth/jwt"
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

func SetAuthToken(ctx context.Context, userId uuid.UUID) error {
	session := GetSession(ctx)
	token, err := jwt.CreateJWT(userId)

	if err != nil {
		return err
	}

	cookie := CreateAuthCookie(token)

	http.SetCookie(session.Writer, &cookie)

	return nil
}

func RemoveAuthToken(ctx context.Context) {
	session := GetSession(ctx)
	cookie := ResetAuthCookie()

	http.SetCookie(session.Writer, &cookie)
}

func GetAuthTokenString(ctx *gin.Context) (string, error) {
	return GetAuthCookieToken(ctx)
}
