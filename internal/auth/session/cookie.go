package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CookieConfig struct {
	Key      string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

const AuthTokenCookieKey = "auth-token"

func GetAuthCookieToken(ctx *gin.Context) (string, error) {
	value, err := ctx.Cookie(AuthTokenCookieKey)

	if err != nil {
		return "", err
	}

	return value, nil
}

func CreateAuthCookie(value string) http.Cookie {
	return http.Cookie{
		Name:     AuthTokenCookieKey,
		Value:    value,
		MaxAge:   0,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
}

func ResetAuthCookie() http.Cookie {
	return http.Cookie{
		Name:     AuthTokenCookieKey,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}
}
