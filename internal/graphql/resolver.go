package graphql

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	"github.com/proctorinc/banker/internal/upload"
)

type Resolver struct {
	Repository  db.Repository
	AuthService auth.AuthService
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

func (r *queryResolver) Transaction(ctx context.Context, transactionId uuid.UUID) (*db.Transaction, error) {
	transaction, err := r.Repository.GetTransaction(ctx, transactionId)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *queryResolver) Transactions(ctx context.Context) ([]db.Transaction, error) {
	user := auth.GetCurrentUser(ctx)

	return r.Repository.ListTransactions(ctx, user.ID)
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Login(ctx context.Context, data LoginInput) (*db.User, error) {
	return r.AuthService.Login(ctx, auth.LoginInput{
		Email:    data.Email,
		Password: data.Password,
	})
}

func (r *mutationResolver) Logout(ctx context.Context) (string, error) {
	return r.AuthService.Logout(ctx)
}

func (r *mutationResolver) Register(ctx context.Context, data RegisterInput) (*db.User, error) {
	return r.AuthService.Register(ctx, auth.RegisterInput{
		Email:    data.Email,
		Username: data.Username,
		Password: data.Password,
	})
}

func (r *mutationResolver) DeleteUser(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	user, err := r.Repository.DeleteUser(ctx, userId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *mutationResolver) CreateTransaction(ctx context.Context, data TransactionInput) (*db.Transaction, error) {
	user := auth.GetCurrentUser(ctx)

	transaction, err := r.Repository.CreateTransaction(ctx, db.CreateTransactionParams{
		Ownerid: user.ID,
		Amount:  int32(data.Amount * 100),
	})

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *mutationResolver) DeleteTransaction(ctx context.Context, id uuid.UUID) (*db.Transaction, error) {
	transaction, err := r.Repository.DeleteTransaction(ctx, id)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *mutationResolver) ChaseTransactionsUpload(ctx context.Context, reader graphql.Upload) (bool, error) {
	user := auth.GetCurrentUser(ctx)
	transactions, err := upload.ParseChaseCSV(reader.File)

	if err != nil {
		return false, err
	}

	for _, transaction := range transactions {
		r.Repository.CreateTransaction(ctx, db.CreateTransactionParams{
			Ownerid: user.ID,
			Amount:  int32(transaction.Amount * 100),
		})
	}

	return true, nil
}

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
	return float64(transaction.Amount) / 100, nil
}
