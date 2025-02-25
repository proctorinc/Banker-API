package resolvers

import (
	"context"
	"database/sql"
	"time"

	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/paging"
	"github.com/proctorinc/banker/internal/graphql/utils"
)

func (r *fundResolver) ID(ctx context.Context, fund *db.Fund) (string, error) {
	return fund.ID.String(), nil
}

func (r *fundResolver) Type(ctx context.Context, fund *db.Fund) (string, error) {
	return string(fund.Type), nil
}

func (r *fundResolver) Name(ctx context.Context, fund *db.Fund) (string, error) {
	return fund.Name, nil
}

func (r *fundResolver) Goal(ctx context.Context, fund *db.Fund) (float64, error) {
	return utils.FormatCurrencyFloat64(fund.Goal * 100), nil
}

func (r *fundResolver) Allocations(ctx context.Context, fund *db.Fund, page *paging.PageArgs) (*gen.FundAllocationConnection, error) {
	totalCount, err := r.DataLoaders.Retrieve(ctx).CountFundAllocationsByFundId.Load(fund.ID.String())

	if err != nil {
		return &gen.FundAllocationConnection{
			PageInfo: paging.NewEmptyPageInfo(),
		}, err
	}

	paginator := paging.NewOffsetPaginator(page, totalCount)
	result := &gen.FundAllocationConnection{
		PageInfo: &paginator.PageInfo,
	}
	start := int32(paginator.Offset)
	limit := calculatePageLimit(page)

	allocations, err := r.DataLoaders.Retrieve(ctx).FundAllocationsByFundId(limit, start).Load(fund.ID.String())

	for i, row := range allocations {
		result.Edges = append(result.Edges, gen.FundAllocationEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	return result, err
}

func (r *fundResolver) StartDate(ctx context.Context, fund *db.Fund) (string, error) {
	return fund.Startdate.Format(time.RFC3339), nil
}

func (r *fundResolver) EndDate(ctx context.Context, fund *db.Fund) (string, error) {
	if fund.Enddate.Valid {
		return fund.Enddate.Time.Format(time.RFC3339), nil
	}

	return "", nil
}

func (r *fundResolver) Total(ctx context.Context, fund *db.Fund) (float64, error) {
	total, err := r.Repository.GetFundTotal(ctx, fund.ID)

	if err != nil {
		return 0, err
	}

	return utils.FormatCurrencyFloat64(int32(total.(int64))), nil
}

// Mutations

func (r *mutationResolver) CreateFund(ctx context.Context, data gen.CreateFundInput) (*db.Fund, error) {
	user := auth.GetCurrentUser(ctx)

	fund, err := r.Repository.CreateFund(ctx, db.CreateFundParams{
		Type:      db.FundType(data.Type),
		Name:      data.Name,
		Goal:      int32(data.Goal),
		Startdate: time.Now(),
		Enddate:   sql.NullTime{Time: time.Now(), Valid: false},
		Ownerid:   user.ID,
	})

	return &fund, err
}
