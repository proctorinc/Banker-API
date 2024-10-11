package jwt

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const JwtSecret = "CHANGE_SECRET"

type JwtToken struct {
	UserId    uuid.UUID
	ExpiresAt time.Time
	IssuedAt  time.Time
}

func VerifyJwt(jwtString string) (*JwtToken, error) {
	jwtToken := new(JwtToken)

	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JwtSecret), nil
	})

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("unable to parse jwt token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId, err := parseUserId(claims)

		if err != nil {
			return nil, fmt.Errorf("unable to parse userId as UUID")
		}

		expiresAt, err := parseExpirationDate(claims)

		if err != nil {
			return nil, fmt.Errorf("unable to parse expiration date")
		}

		issuedAt, err := parseIssuedAt(claims)

		if err != nil {
			return nil, fmt.Errorf("unable to parse expiration date")
		}

		jwtToken = &JwtToken{
			UserId:    *userId,
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt,
		}
	}

	return jwtToken, nil
}

func CreateJWT(userId uuid.UUID) (string, error) {
	log.Println(time.Now().Add(time.Hour * 24).Format(time.RFC3339))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	secret := []byte(JwtSecret)
	return token.SignedString(secret)
}

func parseUserId(claims jwt.MapClaims) (*uuid.UUID, error) {
	idString, err := claims.GetSubject()

	if err != nil {
		return nil, err
	}

	uuid, err := uuid.Parse(idString)

	if err != nil {
		return nil, err
	}

	return &uuid, nil
}

func parseExpirationDate(claims jwt.MapClaims) (time.Time, error) {
	date, err := claims.GetExpirationTime()

	return date.Time, err
}

func parseIssuedAt(claims jwt.MapClaims) (time.Time, error) {
	date, err := claims.GetIssuedAt()

	return date.Time, err
}
