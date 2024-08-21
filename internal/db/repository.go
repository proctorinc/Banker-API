package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository interface {
	// Users
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (User, error)

	// Transactions
	GetTransaction(ctx context.Context, id uuid.UUID) (Transaction, error)
	ListTransactions(ctx context.Context, ownerid uuid.UUID) ([]Transaction, error)
	CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) (Transaction, error)
}

type repositoryService struct {
	*Queries
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repositoryService{
		Queries: New(db),
		db:      db,
	}
}

func Open(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}

func (r repositoryService) withTx(ctx context.Context, txFn func(*Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = txFn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			err = fmt.Errorf("tx failed: %v, unable to rollback: %v", err, rbErr)
		}
	} else {
		err = tx.Commit()
	}
	return err
}
