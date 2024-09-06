package token

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	log.Println("Setting auth token is not implemented")
	// cookie := http.Cookie{
	// 	Name:     TokenKey,
	// 	Value:    userId.String(),
	// 	MaxAge:   3600,
	// 	Path:     "/",
	// 	Domain:   "localhost",
	// 	Secure:   false,
	// 	HttpOnly: true,
	// }

	// http.SetCookie(ctx.Writer, &cookie)

	// ctx.SetCookie(
	// 	config.Key,
	// 	config.Value,
	// 	config.MaxAge,
	// 	config.Path,
	// 	config.Domain,
	// 	config.Secure,
	// 	config.HttpOnly,
	// )
}
