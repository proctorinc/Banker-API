package resolvers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/paging"
	"github.com/proctorinc/banker/internal/graphql/utils"
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

func (r *transactionResolver) Amount(ctx context.Context, transaction *db.Transaction) (float64, error) {
	return utils.FormatCurrencyFloat64(transaction.Amount), nil
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
	merchant, err := r.DataLoaders.Retrieve(ctx).MerchantByTransactionId.Load(transaction.Merchantid.String())

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

func (r *queryResolver) Transactions(ctx context.Context, page *paging.PageArgs) (*gen.TransactionConnection, error) {
	user := auth.GetCurrentUser(ctx)
	totalCount, err := r.Repository.CountTransactions(ctx, user.ID)

	if err != nil {
		return &gen.TransactionConnection{
			PageInfo: paging.NewEmptyPageInfo(),
		}, err
	}

	paginator := paging.NewOffsetPaginator(page, totalCount)
	result := &gen.TransactionConnection{
		PageInfo: &paginator.PageInfo,
	}
	start := int32(paginator.Offset)
	limit := calculatePageLimit(page)

	transactions, err := r.Repository.ListTransactions(ctx, db.ListTransactionsParams{
		Ownerid: user.ID,
		Limit:   limit,
		Start:   start,
	})

	for i, row := range transactions {
		result.Edges = append(result.Edges, gen.TransactionEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	return result, err
}

// Mutations

func (r *mutationResolver) DeleteTransaction(ctx context.Context, id uuid.UUID) (*db.Transaction, error) {
	transaction, err := r.Repository.DeleteTransaction(ctx, id)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
