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

	// Accounts
	GetAccount(ctx context.Context, arg GetAccountParams) (Account, error)
	ListAccounts(ctx context.Context, ownerid uuid.UUID) ([]Account, error)
	UpsertAccount(ctx context.Context, arg UpsertAccountParams) (Account, error)
	GetAccountSpending(ctx context.Context, arg GetAccountSpendingParams) (interface{}, error) //(GetAccountSpendingRow, error)
	GetAccountIncome(ctx context.Context, arg GetAccountIncomeParams) (interface{}, error)     //(GetAccountIncomeRow, error)

	// Transactions
	GetTransaction(ctx context.Context, arg GetTransactionParams) (Transaction, error)
	ListTransactions(ctx context.Context, ownerid uuid.UUID) ([]Transaction, error)
	UpsertTransaction(ctx context.Context, arg UpsertTransactionParams) (Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) (Transaction, error)
	GetTotalSpending(ctx context.Context, id uuid.UUID) (interface{}, error) //(GetTotalSpendingRow, error)
	GetTotalIncome(ctx context.Context, id uuid.UUID) (interface{}, error)   //(GetTotalIncomeRow, error)

	// Merchants
	GetMerchant(ctx context.Context, arg GetMerchantParams) (Merchant, error)
	GetMerchantByKey(ctx context.Context, arg GetMerchantByKeyParams) (Merchant, error)
	ListMerchants(ctx context.Context, ownerid uuid.UUID) ([]Merchant, error)
	CreateMerchant(ctx context.Context, arg CreateMerchantParams) (Merchant, error)

	// Merchant keys
	CreateMerchantKey(ctx context.Context, arg CreateMerchantKeyParams) (MerchantKey, error)
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

func (r *repositoryService) LinkMerchant(ctx context.Context) (*Merchant, error) {
	book := new(Book)
	err := r.withTx(ctx, func(q *Queries) error {
		res, err := q.CreateBook(ctx, bookArg)
		if err != nil {
			return err
		}
		for _, authorID := range authorIDs {
			if err := q.SetBookAuthor(ctx, SetBookAuthorParams{
				BookID:   res.ID,
				AuthorID: authorID,
			}); err != nil {
				return err
			}
		}
		book = &res
		return nil
	})
	return book, err
}
