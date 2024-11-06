package dataloaders

//go:generate go run github.com/vektah/dataloaden TransactionLoader string []github.com/proctorinc/banker/internal/db.Transaction
//go:generate go run github.com/vektah/dataloaden TransactionCountLoader string int64

import (
	"context"
	"time"

	"github.com/proctorinc/banker/internal/db"
)

type contextKey string

const key = contextKey("dataloaders")

// Loaders holds references to the individual dataloaders.
type Loaders struct {
	TransactionsByAccountId      func(limit int32, start int32) *TransactionLoader
	CountTransactionsByAccountId *TransactionCountLoader
}

func newLoaders(ctx context.Context, repo db.Repository) *Loaders {
	return &Loaders{
		TransactionsByAccountId: func(limit int32, start int32) *TransactionLoader {
			return newTransactionsByAccountIdLoader(ctx, repo, limit, start)
		},
		CountTransactionsByAccountId: newCountTransactionsByAccountIdLoader(ctx, repo),
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
				Limit:      limit,
				Start:      start,
			})

			if err != nil {
				return nil, []error{err}
			}

			groupByAccountId := make(map[string][]db.Transaction, len(accountIds))

			for _, r := range res {
				groupByAccountId[r.Accountid.String()] = append(groupByAccountId[r.Accountid.String()], r)
			}

			result := make([][]db.Transaction, len(accountIds))

			for i, accountId := range accountIds {
				result[i] = groupByAccountId[accountId]
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
