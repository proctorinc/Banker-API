package dataloaders

//go:generate go run github.com/vektah/dataloaden TransactionLoader string []github.com/proctorinc/banker/internal/db.Transaction
//go:generate go run github.com/vektah/dataloaden TransactionCountLoader string int64
//go:generate go run github.com/vektah/dataloaden MerchantLoader string github.com/proctorinc/banker/internal/db.Merchant

import (
	"context"
	"time"

	"github.com/proctorinc/banker/internal/db"
)

type contextKey string

const key = contextKey("dataloaders")

// Loaders holds references to the individual dataloaders.
type Loaders struct {
	TransactionsByAccountId       func(limit int32, start int32) *TransactionLoader
	TransactionsByMerchantId      func(limit int32, start int32) *TransactionLoader
	CountTransactionsByAccountId  *TransactionCountLoader
	CountTransactionsByMerchantId *TransactionCountLoader
	MerchantByTransactionId       *MerchantLoader
}

func newLoaders(ctx context.Context, repo db.Repository) *Loaders {
	return &Loaders{
		TransactionsByAccountId: func(limit int32, start int32) *TransactionLoader {
			return newTransactionsByAccountIdLoader(ctx, repo, limit, start)
		},
		TransactionsByMerchantId: func(limit int32, start int32) *TransactionLoader {
			return newTransactionsByMerchantIdLoader(ctx, repo, limit, start)
		},
		CountTransactionsByAccountId:  newCountTransactionsByAccountIdLoader(ctx, repo),
		CountTransactionsByMerchantId: newCountTransactionsByMerchantIdLoader(ctx, repo),
		MerchantByTransactionId:       newMerchantLoader(ctx, repo),
	}
}

type Retriever interface {
	Retrieve(context.Context) *Loaders
}

type retriever struct {
	key contextKey
}

func (r *retriever) Retrieve(ctx context.Context) *Loaders {
	return ctx.Value(r.key).(*Loaders)
}

func NewRetriever() Retriever {
	return &retriever{key: key}
}

func newTransactionsByAccountIdLoader(ctx context.Context, repo db.Repository, limit int32, start int32) *TransactionLoader {
	return NewTransactionLoader(TransactionLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(accountIds []string) ([][]db.Transaction, []error) {
			res, err := repo.ListTransactionsByAccountIds(ctx, db.ListTransactionsByAccountIdsParams{
				Accountids: accountIds,
				Limit:      limit * int32(len(accountIds)), // Query for up to the combined limit
				Start:      start,
			})

			if err != nil {
				return nil, []error{err}
			}

			groupByAccountId := make(map[string][]db.Transaction, len(accountIds))

			for _, r := range res {
				transactions := groupByAccountId[r.Accountid.String()]

				// Make sure transactions are only added up to the limit
				if len(transactions) < int(limit) {
					groupByAccountId[r.Accountid.String()] = append(groupByAccountId[r.Accountid.String()], r)
				}
			}

			result := make([][]db.Transaction, len(accountIds))

			for i, accountId := range accountIds {
				result[i] = groupByAccountId[accountId]
			}

			return result, nil
		},
	})
}

func newTransactionsByMerchantIdLoader(ctx context.Context, repo db.Repository, limit int32, start int32) *TransactionLoader {
	return NewTransactionLoader(TransactionLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(merchantIds []string) ([][]db.Transaction, []error) {
			res, err := repo.ListTransactionsByMerchantIds(ctx, db.ListTransactionsByMerchantIdsParams{
				Merchantids: merchantIds,
				Limit:       limit,
				Start:       start,
			})

			if err != nil {
				return nil, []error{err}
			}

			groupByMerchantId := make(map[string][]db.Transaction, len(merchantIds))

			for _, r := range res {
				groupByMerchantId[r.Merchantid.String()] = append(groupByMerchantId[r.Merchantid.String()], r)
			}

			result := make([][]db.Transaction, len(merchantIds))

			for i, merchantId := range merchantIds {
				result[i] = groupByMerchantId[merchantId]
			}

			return result, nil
		},
	})
}

func newCountTransactionsByAccountIdLoader(ctx context.Context, repo db.Repository) *TransactionCountLoader {
	return NewTransactionCountLoader(TransactionCountLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(accountIds []string) ([]int64, []error) {
			res, err := repo.CountTransactionsByAccountIds(ctx, accountIds)

			if err != nil {
				return nil, []error{err}
			}

			counts := make([]int64, len(accountIds))

			for i, r := range res {
				counts[i] = r.Count
			}

			return counts, nil
		},
	})
}

func newCountTransactionsByMerchantIdLoader(ctx context.Context, repo db.Repository) *TransactionCountLoader {
	return NewTransactionCountLoader(TransactionCountLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(merchantIds []string) ([]int64, []error) {
			res, err := repo.CountTransactionsByMerchantIds(ctx, merchantIds)

			if err != nil {
				return nil, []error{err}
			}

			counts := make([]int64, len(merchantIds))

			for i, r := range res {
				counts[i] = r.Count
			}

			return counts, nil
		},
	})
}

func newMerchantLoader(ctx context.Context, repo db.Repository) *MerchantLoader {
	return NewMerchantLoader(MerchantLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(merchantIds []string) ([]db.Merchant, []error) {
			res, err := repo.ListMerchantsByMerchantIds(ctx, merchantIds)

			if err != nil {
				return nil, []error{err}
			}

			groupByMerchantId := make(map[string]db.Merchant, len(merchantIds))

			for i, r := range res {
				groupByMerchantId[r.ID.String()] = res[i]
			}

			result := make([]db.Merchant, len(merchantIds))

			for i, merchantId := range merchantIds {
				result[i] = groupByMerchantId[merchantId]
			}

			return result, nil
		},
	})
}
