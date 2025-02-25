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
	UpsertAccount(ctx context.Context, arg UpsertAccountParams) (Account, error)
	CountAccounts(ctx context.Context, ownerid uuid.UUID) (int64, error)

	// Account Sync Item
	GetLastSync(ctx context.Context, accountId uuid.UUID) (AccountSyncItem, error)
	CreateAccountSyncItem(ctx context.Context, arg CreateAccountSyncItemParams) (AccountSyncItem, error)

	// Transactions
	GetTransaction(ctx context.Context, arg GetTransactionParams) (Transaction, error)
	ListTransactions(ctx context.Context, arg ListTransactionsParams) ([]Transaction, error)
	ListTransactionsByDates(ctx context.Context, arg ListTransactionsByDatesParams) ([]Transaction, error)
	ListTransactionsByAccountIds(ctx context.Context, arg ListTransactionsByAccountIdsParams) ([]Transaction, error)
	ListTransactionsByMerchantIds(ctx context.Context, arg ListTransactionsByMerchantIdsParams) ([]Transaction, error)
	ListSpendingTransactions(ctx context.Context, arg ListSpendingTransactionsParams) ([]Transaction, error)
	ListIncomeTransactions(ctx context.Context, arg ListIncomeTransactionsParams) ([]Transaction, error)
	ListAccountSpendingTransactions(ctx context.Context, arg ListAccountSpendingTransactionsParams) ([]Transaction, error)
	ListAccountIncomeTransactions(ctx context.Context, arg ListAccountIncomeTransactionsParams) ([]Transaction, error)
	ListMonths(ctx context.Context, args uuid.UUID) ([]ListMonthsRow, error)
	CountTransactions(ctx context.Context, ownerid uuid.UUID) (int64, error)
	CountTransactionsByDates(ctx context.Context, arg CountTransactionsByDatesParams) (int64, error)
	CountTransactionsByAccountIds(ctx context.Context, accountIds []string) ([]CountTransactionsByAccountIdsRow, error)
	CountTransactionsByMerchantIds(ctx context.Context, merchantIds []string) ([]CountTransactionsByMerchantIdsRow, error)
	CountIncomeTransactions(ctx context.Context, arg CountIncomeTransactionsParams) (int64, error)
	CountSpendingTransactions(ctx context.Context, arg CountSpendingTransactionsParams) (int64, error)
	UpsertTransaction(ctx context.Context, arg UpsertTransactionParams) (Transaction, error)
	DeleteTransaction(ctx context.Context, id uuid.UUID) (Transaction, error)

	// Merchants
	GetMerchant(ctx context.Context, arg GetMerchantParams) (Merchant, error)
	GetMerchantByName(ctx context.Context, name string) (Merchant, error)
	GetMerchantBySourceId(ctx context.Context, sourceId sql.NullString) (Merchant, error)
	GetMerchantByKey(ctx context.Context, arg GetMerchantByKeyParams) (Merchant, error)
	ListMerchants(ctx context.Context, arg ListMerchantsParams) ([]Merchant, error)
	ListMerchantsByMerchantIds(ctx context.Context, merchantIds []string) ([]Merchant, error)
	CountMerchants(ctx context.Context, ownerid uuid.UUID) (int64, error)
	CreateMerchant(ctx context.Context, arg CreateMerchantParams) (Merchant, error)
	LinkMerchant(ctx context.Context, arg LinkMerchantParams) (*Merchant, error)

	// Merchant keys
	CreateMerchantKey(ctx context.Context, arg CreateMerchantKeyParams) (MerchantKey, error)

	// Stats
	GetTotalSpending(ctx context.Context, arg GetTotalSpendingParams) (interface{}, error)
	GetTotalIncome(ctx context.Context, arg GetTotalIncomeParams) (interface{}, error)
	GetNetIncome(ctx context.Context, arg GetNetIncomeParams) (interface{}, error)
	GetAccountSpending(ctx context.Context, arg GetAccountSpendingParams) (interface{}, error)
	GetAccountIncome(ctx context.Context, arg GetAccountIncomeParams) (interface{}, error)

	// Funds
	CreateFund(ctx context.Context, arg CreateFundParams) (Fund, error)
	ListSavingsFunds(ctx context.Context, arg ListSavingsFundsParams) ([]Fund, error)
	ListBudgetFunds(ctx context.Context, arg ListBudgetFundsParams) ([]Fund, error)
	GetFundTotal(ctx context.Context, fundId uuid.UUID) (interface{}, error)
	CountSavingsFunds(ctx context.Context, ownerid uuid.UUID) (int64, error)
	CountBudgetFunds(ctx context.Context, ownerid uuid.UUID) (int64, error)

	// Fund Allocations
	ListFundAllocationsByFundIds(ctx context.Context, arg ListFundAllocationsByFundIdsParams) ([]FundAllocation, error)
	CountFundAllocationsByFundId(ctx context.Context, fundIds []string) ([]CountFundAllocationsByFundIdRow, error)
	GetFundAllocationsStats(ctx context.Context, arg GetFundAllocationsStatsParams) (GetFundAllocationsStatsRow, error)
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
