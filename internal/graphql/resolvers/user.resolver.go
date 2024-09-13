package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
)

type userResolver struct{ *Resolver }

func (r *userResolver) ID(ctx context.Context, user *db.User) (string, error) {
	return user.ID.String(), nil
}

func (r *userResolver) Role(ctx context.Context, user *db.User) (string, error) {
	if user.Role == db.RoleADMIN {
		return "Admin", nil
	}

	return "User", nil
}

func (r *userResolver) Accounts(ctx context.Context, user *db.User) ([]db.Account, error) {
	transactions, err := r.Repository.ListAccounts(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *userResolver) Transactions(ctx context.Context, user *db.User) ([]db.Transaction, error) {
	transactions, err := r.Repository.ListTransactions(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// Queries

func (r *queryResolver) Me(ctx context.Context) (*db.User, error) {
	user := auth.GetCurrentUser(ctx)

	if user == nil {
		return nil, fmt.Errorf("no authenticated user")
	}

	return user, nil
}

func (r *queryResolver) User(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	user, err := r.Repository.GetUser(ctx, userId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]db.User, error) {
	return r.Repository.ListUsers(ctx)
}

// Mutations

func (r *mutationResolver) Login(ctx context.Context, data gen.LoginInput) (*db.User, error) {
	return r.AuthService.Login(ctx, auth.LoginInput{
		Email:    data.Email,
		Password: data.Password,
	})
}

func (r *mutationResolver) Logout(ctx context.Context) (string, error) {
	return r.AuthService.Logout(ctx)
}

func (r *mutationResolver) Register(ctx context.Context, data gen.RegisterInput) (*db.User, error) {
	return r.AuthService.Register(ctx, auth.RegisterInput{
		Email:    data.Email,
		Username: data.Username,
		Password: data.Password,
	})
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (*db.User, error) {
	user := auth.GetCurrentUser(ctx)
	deleted, err := r.Repository.DeleteUser(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	return &deleted, nil
}
