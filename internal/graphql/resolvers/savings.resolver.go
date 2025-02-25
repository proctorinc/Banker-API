package resolvers

import (
	"context"

	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/paging"
)

func (r *queryResolver) SavingsFunds(ctx context.Context, filter gen.DateFilter) (*gen.FundsResponse, error) {
	user := auth.GetCurrentUser(ctx)
	result := gen.FundsResponse{}

	dateRange, err := parseStatsFilter(&filter)

	if err != nil {
		return nil, err
	}

	stats, err := r.Repository.GetFundAllocationsStats(ctx, db.GetFundAllocationsStatsParams{
		Ownerid: user.ID,
		Enddate: dateRange.EndDate,
	})

	if err != nil {
		return nil, err
	}

	result.Stats = &gen.FundsStats{
		TotalSavings: float64(stats.Net.(int64)),
		Saved:        float64(stats.Saved.(int64)),
		Spent:        float64(stats.Spent.(int64)),
		Unallocated:  0,
	}
	// Funds to be resolver by fundsResponseResolver below
	result.Funds = &gen.FundConnection{}

	return &result, nil
}

func (r fundsResponseResolver) Funds(ctx context.Context, response *gen.FundsResponse, page *paging.PageArgs) (*gen.FundConnection, error) {
	user := auth.GetCurrentUser(ctx)
	totalCount, err := r.Repository.CountSavingsFunds(ctx, user.ID)

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

	transactions, err := r.Repository.ListSavingsFunds(ctx, db.ListSavingsFundsParams{
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
