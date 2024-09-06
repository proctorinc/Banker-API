package token

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth/session"
)

type AuthToken struct {
	Value string
}

type CookieConfig struct {
	Key      string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

const TokenKey = "auth-token"

func (t *AuthToken) IsEmpty() bool {
	return t != nil && t.Value == ""
}

func (t *AuthToken) GetUserId() (uuid.UUID, error) {
	return uuid.Parse(t.Value)
}

func GetAuthToken(ctx *gin.Context) (*AuthToken, error) {
	value, err := ctx.Cookie(TokenKey)

	if err != nil {
		return nil, err
	}

	token := &AuthToken{
		Value: value,
	}

	return token, nil
}

func SetAuthToken(ctx context.Context, userId uuid.UUID) {
	session := session.GetSession(ctx)

	cookie := http.Cookie{
		Name:     TokenKey,
		Value:    userId.String(),
		MaxAge:   0,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(session.Writer, &cookie)
}

func RemoveAuthToken(ctx context.Context) {
	session := session.GetSession(ctx)

	cookie := http.Cookie{
		Name:     TokenKey,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(session.Writer, &cookie)
}
