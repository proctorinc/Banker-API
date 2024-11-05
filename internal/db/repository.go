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
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (User, error)

	// Accounts
	GetAccount(ctx context.Context, arg GetAccountParams) (Account, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListAccountsCount(ctx context.Context, ownerid uuid.UUID) (int64, error)
	UpsertAccount(ctx context.Context, arg UpsertAccountParams) (Account, error)
	GetAccountSpending(ctx context.Context, arg GetAccountSpendingParams) (interface{}, error)
	GetAccountIncome(ctx context.Context, arg GetAccountIncomeParams) (interface{}, error)

	// Transactions
	GetTransaction(ctx context.Context, arg GetTransactionParams) (Transaction, error)
	ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]Transaction, error)
	ListTransactionsCount(ctx context.Context, ownerid uuid.UUID) (int64, error)
	ListTransactionsByAccountIds(ctx context.Context, arg ListTransactionsByAccountIdsParams) ([]Transaction, error)
	ListTransactionsByMerchantId(ctx context.Context, arg ListTransactionsByMerchantIdParams) ([]Transaction, error)
	ListSpendingTransactions(ctx context.Context, arg ListSpendingTransactionsParams) ([]Transaction, error)
	ListIncomeTransactions(ctx context.Context, arg ListIncomeTransactionsParams) ([]Transaction, error)
	ListAccountSpendingTransactions(ctx context.Context, arg ListAccountSpendingTransactionsParams) ([]Transaction, error)
	ListAccountIncomeTransactions(ctx context.Context, arg ListAccountIncomeTransactionsParams) ([]Transaction, error)
	UpsertTransaction(ctx context.Context, arg UpsertTransactionParams) (Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) (Transaction, error)
	GetTotalSpending(ctx context.Context, arg GetTotalSpendingParams) (interface{}, error)
	GetTotalIncome(ctx context.Context, arg GetTotalIncomeParams) (interface{}, error)

	// Merchants
	GetMerchant(ctx context.Context, arg GetMerchantParams) (Merchant, error)
	GetMerchantByName(ctx context.Context, name string) (Merchant, error)
	GetMerchantBySourceId(ctx context.Context, sourceId sql.NullString) (Merchant, error)
	GetMerchantByKey(ctx context.Context, arg GetMerchantByKeyParams) (Merchant, error)
	ListMerchants(ctx context.Context, arg ListMerchantsParams) ([]Merchant, error)
	ListMerchantsCount(ctx context.Context, ownerid uuid.UUID) (int64, error)
	CreateMerchant(ctx context.Context, arg CreateMerchantParams) (Merchant, error)
	LinkMerchant(ctx context.Context, arg LinkMerchantParams) (*Merchant, error)

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

type LinkMerchantParams struct {
	MerchantName string
	KeyMatch     string
	UploadSource UploadSource
	SourceId     sql.NullString
	UserId       uuid.UUID
}

func (r *repositoryService) LinkMerchant(ctx context.Context, arg LinkMerchantParams) (*Merchant, error) {
	merchant := new(Merchant)

	err := r.withTx(ctx, func(q *Queries) error {
		res, err := q.CreateMerchant(ctx, CreateMerchantParams{
			Name:     arg.MerchantName,
			Ownerid:  arg.UserId,
			Sourceid: arg.SourceId,
		})

		if err != nil {
			return err
		}

		_, err = q.CreateMerchantKey(ctx, CreateMerchantKeyParams{
			Keymatch:     arg.KeyMatch,
			Uploadsource: arg.UploadSource,
			Merchantid:   res.ID,
			Ownerid:      arg.UserId,
		})

		if err != nil {
			return err
		}
		merchant = &res
		return nil
	})
	return merchant, err
}
