package graphql

import (
	"context"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/db"
)

type Resolver struct {
	Repository db.Repository
	// DataLoaders dataloaders.Retriever
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

func (r *Resolver) Transaction() TransactionResolver {
	return &transactionResolver{r}
}

type queryResolver struct{ *Resolver }

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

func (r *queryResolver) Transaction(ctx context.Context, transactionId uuid.UUID) (*db.Transaction, error) {
	transaction, err := r.Repository.GetTransaction(ctx, transactionId)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *queryResolver) Transactions(ctx context.Context, userId uuid.UUID) ([]db.Transaction, error) {
	return r.Repository.ListTransactions(ctx, userId)
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Login(ctx context.Context, data LoginInput) (*db.User, error) {
	user, err := r.Repository.GetUserByEmail(ctx, data.Email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

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

func (r *mutationResolver) DeleteUser(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	user, err := r.Repository.DeleteUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mutationResolver) CreateTransaction(ctx context.Context, data TransactionInput) (*db.Transaction, error) {
	user, err := r.Repository.CreateTransaction(ctx, db.CreateTransactionParams{
		Amount:  int32(data.Amount * 100),
		Ownerid: data.OwnerID,
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *mutationResolver) DeleteTransaction(ctx context.Context, id uuid.UUID) (*db.Transaction, error) {
	transaction, err := r.Repository.DeleteTransaction(ctx, id)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) ID(ctx context.Context, user *db.User) (string, error) {
	return user.ID.String(), nil
}

func (r *userResolver) Transactions(ctx context.Context, user *db.User) ([]db.Transaction, error) {
	transactions, err := r.Repository.ListTransactions(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

type transactionResolver struct{ *Resolver }

func (r *transactionResolver) ID(ctx context.Context, transaction *db.Transaction) (uuid.UUID, error) {
	return transaction.ID, nil
}

func (r *transactionResolver) Amount(ctx context.Context, transaction *db.Transaction) (float64, error) {
	return float64(transaction.Amount / 100), nil
}
