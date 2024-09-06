package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/auth/password"
	"github.com/proctorinc/banker/internal/auth/token"
	"github.com/proctorinc/banker/internal/db"
)

type AuthService struct {
	Repository db.Repository
}

type LoginInput struct {
	Email    string
	Password string
}

type RegisterInput struct {
	Email    string
	Username string
	Password string
}

const userContextKey = "user"

func NewAuthService(r db.Repository) *AuthService {
	return &AuthService{
		Repository: r,
	}
}

func IsAuthenticated(ctx context.Context) bool {
	return GetAuthenticatedUser(ctx) != nil
}

func GetAuthenticatedUser(ctx context.Context) *db.User {
	raw, _ := ctx.Value(userContextKey).(*db.User)
	return raw
}

func SetAuthenticatedUser(ctx *gin.Context, user db.User) {
	ctx.Set(userContextKey, user)
}

func (s *AuthService) GetCurrentUser(ctx context.Context) *db.User {
	return GetAuthenticatedUser(ctx)
}

func (s *AuthService) Login(ctx context.Context, data LoginInput) (*db.User, error) {
	user, err := s.Repository.GetUserByEmail(ctx, data.Email)

	if err != nil {
		log.Printf("Invalid email login from %s", data.Email)
		return nil, fmt.Errorf("Login failed. Invalid email or password")
	}

	if err = password.VerifyPassword(data.Password, user.Passwordhash); err != nil {
		log.Printf("Login failed. Invalid password for user: %s", user.Email)
		return nil, fmt.Errorf("Login failed. Invalid email or password")
	}

	token.SetAuthToken(ctx, user.ID)

	return &user, nil
}

func (s *AuthService) Logout(ctx context.Context) (string, error) {
	return "User has been successfully logged out", nil
}

func (s *AuthService) Register(ctx context.Context, data RegisterInput) (*db.User, error) {
	hash, err := password.HashPassword(data.Password)

	if err != nil {
		return nil, err
	}

	params := db.CreateUserParams{
		Email:        data.Email,
		Username:     data.Username,
		Passwordhash: hash,
	}

	user, err := s.Repository.CreateUser(ctx, params)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
