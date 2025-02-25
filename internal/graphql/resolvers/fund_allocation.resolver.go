package resolvers

import (
	"context"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/paging"
	"github.com/proctorinc/banker/internal/graphql/utils"
)

func (r *fundAllocationResolver) ID(ctx context.Context, allocation *db.FundAllocation) (string, error) {
	return allocation.ID.String(), nil
}

func (r *fundAllocationResolver) Description(ctx context.Context, allocation *db.FundAllocation) (string, error) {
	return allocation.Description, nil
}

func (r *fundAllocationResolver) Amount(ctx context.Context, allocation *db.FundAllocation) (float64, error) {
	return utils.FormatCurrencyFloat64(allocation.Amount), nil
}

// Queries
func (r *queryResolver) Fund(ctx context.Context, fundId uuid.UUID) (*db.Fund, error) {
	return &db.Fund{}, nil
}

// func (r *queryResolver) SavingsFunds(ctx context.Context, page *paging.PageArgs) (*gen.FundConnection, error) {
// 	user := auth.GetCurrentUser(ctx)
// 	totalCount, err := r.Repository.CountSavingsFunds(ctx, user.ID)

// 	if err != nil {
// 		return &gen.FundConnection{
// 			PageInfo: paging.NewEmptyPageInfo(),
// 		}, err
// 	}

// 	paginator := paging.NewOffsetPaginator(page, totalCount)
// 	result := &gen.FundConnection{
// 		PageInfo: &paginator.PageInfo,
// 	}
// 	start := int32(paginator.Offset)
// 	limit := calculatePageLimit(page)

// 	transactions, err := r.Repository.ListSavingsFunds(ctx, db.ListSavingsFundsParams{
// 		Ownerid: user.ID,
// 		Limit:   limit,
// 		Start:   start,
// 	})

// 	for i, row := range transactions {
// 		result.Edges = append(result.Edges, gen.FundEdge{
// 			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
// 			Node:   &row,
// 		})
// 	}

// 	return result, err
// }

func (r *queryResolver) Budgets(ctx context.Context, page *paging.PageArgs) (*gen.FundConnection, error) {
	user := auth.GetCurrentUser(ctx)
	totalCount, err := r.Repository.CountBudgetFunds(ctx, user.ID)

	if err != nil {
		return &gen.FundConnection{
			PageInfo: paging.NewEmptyPageInfo(),
		}, err
	}

	paginator := paging.NewOffsetPaginator(page, totalCount)
	result := &gen.FundConnection{
		PageInfo: &paginator.PageInfo,
	}
	start := int32(paginator.Offset)
	limit := calculatePageLimit(page)

	transactions, err := r.Repository.ListBudgetFunds(ctx, db.ListBudgetFundsParams{
		Ownerid: user.ID,
		Limit:   limit,
		Start:   start,
	})

	for i, row := range transactions {
		result.Edges = append(result.Edges, gen.FundEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	return result, err
}
