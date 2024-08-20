package graphql

import (
	"context"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/db"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository db.Repository
	// DataLoaders dataloaders.Retriever
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

func (r *queryResolver) User(ctx context.Context, id string) (*db.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	user, err := r.Repository.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]db.User, error) {
	return r.Repository.ListUsers(ctx)
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, data UserInput) (*db.User, error) {
	user, err := r.Repository.CreateUser(ctx, db.CreateUserParams{
		Username: data.Username,
		Email:    data.Email,
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*db.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	user, err := r.Repository.DeleteUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) ID(ctx context.Context, obj *db.User) (string, error) {
	return obj.ID.String(), nil
}
