package resolvers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
)

type UploadResponse struct {
	Success              bool
	AccountsUploaded     int
	TransactionsUploaded int
}

var UploadFailed = UploadResponse{
	Success:              false,
	AccountsUploaded:     0,
	TransactionsUploaded: 0,
}

func (r *transactionResolver) ID(ctx context.Context, transaction *db.Transaction) (uuid.UUID, error) {
	return transaction.ID, nil
}

func (r *transactionResolver) SourceId(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Sourceid, nil
}

func (r *transactionResolver) UploadSource(ctx context.Context, transaction *db.Transaction) (string, error) {
	return string(transaction.Uploadsource), nil
}

func (r *transactionResolver) Amount(ctx context.Context, transaction *db.Transaction) (float64, error) {
	return float64(transaction.Amount) / 100, nil
}

func (r *transactionResolver) PayeeID(ctx context.Context, transaction *db.Transaction) (*string, error) {
	if len(transaction.Payeeid.String) > 0 {
		return &transaction.Payeeid.String, nil
	}

	return nil, nil
}

func (r *transactionResolver) Payee(ctx context.Context, transaction *db.Transaction) (*string, error) {
	if len(transaction.Payee.String) > 0 {
		return &transaction.Payee.String, nil
	}

	return nil, nil
}

func (r *transactionResolver) PayeeFull(ctx context.Context, transaction *db.Transaction) (*string, error) {
	if len(transaction.Payeefull.String) > 0 {
		return &transaction.Payeefull.String, nil
	}

	return nil, nil
}

func (r *transactionResolver) IsoCurrencyCode(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Isocurrencycode, nil
}

func (r *transactionResolver) Date(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Date.Format(time.RFC3339), nil
}

func (r *transactionResolver) Description(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Description, nil
}

func (r *transactionResolver) Type(ctx context.Context, transaction *db.Transaction) (string, error) {
	return string(transaction.Type), nil
}

func (r *transactionResolver) CheckNumber(ctx context.Context, transaction *db.Transaction) (*string, error) {
	if len(transaction.Checknumber.String) > 0 {
		return &transaction.Checknumber.String, nil
	}

	return nil, nil
}

func (r *transactionResolver) Updated(ctx context.Context, transaction *db.Transaction) (string, error) {
	return transaction.Updated.Format(time.RFC3339), nil
}

func (r *transactionResolver) Merchant(ctx context.Context, transaction *db.Transaction) (*db.Merchant, error) {
	user := auth.GetCurrentUser(ctx)
	merchant, err := r.Repository.GetMerchant(ctx, db.GetMerchantParams{
		ID:      transaction.Merchantid,
		Ownerid: user.ID,
	})

	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

// Queries

func (r *queryResolver) Transaction(ctx context.Context, transactionId uuid.UUID) (*db.Transaction, error) {
	user := auth.GetCurrentUser(ctx)
	transaction, err := r.Repository.GetTransaction(ctx, db.GetTransactionParams{
		ID:      transactionId,
		Ownerid: user.ID,
	})

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *queryResolver) Transactions(ctx context.Context) ([]db.Transaction, error) {
	user := auth.GetCurrentUser(ctx)

	return r.Repository.ListTransactions(ctx, user.ID)
}

func (r *queryResolver) SpendingTotal(ctx context.Context) (float64, error) {
	user := auth.GetCurrentUser(ctx)

	spending, err := r.Repository.GetTotalSpending(ctx, user.ID)

	if err != nil {
		return 0, err
	}

	return float64(spending.(int64)) / 100, nil
}

func (r *queryResolver) IncomeTotal(ctx context.Context) (float64, error) {
	user := auth.GetCurrentUser(ctx)

	income, err := r.Repository.GetTotalIncome(ctx, user.ID)

	if err != nil {
		return 0, err
	}

	return float64(income.(int64)) / 100, nil
}

func (r *queryResolver) Stats(ctx context.Context, input *gen.StatsInput) (*gen.StatsResponse, error) {
	user := auth.GetCurrentUser(ctx)
	incomeTotal, err := r.Repository.GetTotalIncome(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	income := gen.IncomeStats{
		Total:        float64(incomeTotal.(int64)) / 100,
		Transactions: []db.Transaction{},
	}

	spendingTotal, err := r.Repository.GetTotalSpending(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	spending := gen.SpendingStats{
		Total:        float64(spendingTotal.(int64)) / 100,
		Transactions: []db.Transaction{},
	}

	net := gen.NetStats{
		Total: float64(incomeTotal.(int64)+spendingTotal.(int64)) / 100,
	}

	response := &gen.StatsResponse{
		Spending: &spending,
		Income:   &income,
		Net:      &net,
	}

	return response, nil
}

// Mutations

func (r *mutationResolver) DeleteTransaction(ctx context.Context, id uuid.UUID) (*db.Transaction, error) {
	transaction, err := r.Repository.DeleteTransaction(ctx, id)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
