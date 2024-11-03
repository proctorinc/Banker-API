package dataloaders

//go:generate go run github.com/vektah/dataloaden TransactionLoader string []github.com/proctorinc/banker/internal/db.Transaction
/*
*
* TODO: Add other dataloaders
*
 */

import (
	"context"
	"time"

	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
)

type contextKey string

const key = contextKey("dataloaders")

// Loaders holds references to the individual dataloaders.
type Loaders struct {
	TransactionsByAccountId *TransactionLoader
}

func newLoaders(ctx context.Context, repo db.Repository) *Loaders {
	return &Loaders{
		TransactionsByAccountId: newTransactionsByAccountIdLoader(ctx, repo),
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

// NewRetriever instantiates a new implementation of Retriever.
func NewRetriever() Retriever {
	return &retriever{key: key}
}

func newTransactionsByAccountIdLoader(ctx context.Context, repo db.Repository) *TransactionLoader {
	return NewTransactionLoader(TransactionLoaderConfig{
		MaxBatch: 100,
		Wait:     5 * time.Millisecond,
		Fetch: func(accountIds []string) ([][]db.Transaction, []error) {
			user := auth.GetCurrentUser(ctx)

			// db query
			res, err := repo.ListTransactionsByAccountIds(ctx, db.ListTransactionsByAccountIdsParams{
				Ownerid: user.ID,
				Column2: accountIds,
			})
			if err != nil {
				return nil, []error{err}
			}
			// map
			groupByAccountId := make(map[string][]db.Transaction, len(accountIds))
			for _, r := range res {
				groupByAccountId[r.Accountid.String()] = append(groupByAccountId[r.Accountid.String()], r)
			}
			// order
			result := make([][]db.Transaction, len(accountIds))
			for i, accountId := range accountIds {
				result[i] = groupByAccountId[accountId]
			}
			return result, nil
		},
	})
}
