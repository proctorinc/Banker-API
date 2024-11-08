package resolvers

import (
	"context"

	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/paging"
	"github.com/proctorinc/banker/internal/graphql/utils"
)

type StatsResolver struct {
	Spending gen.SpendingStats
	Income   gen.IncomeStats
	Net      gen.NetStats
}

type FakePageArgs struct {
	First int
	After *string
}

type Page struct {
	First int
}

// Queries

func (r *queryResolver) Spending(ctx context.Context, input gen.StatsInput) (*gen.SpendingStats, error) {
	user := auth.GetCurrentUser(ctx)
	pageArgs := getPageArgs(ctx, "transactions")
	filter, err := parseStatsFilter(input.Filter)

	if err != nil {
		return nil, err
	}

	incomeTotal, err := r.Repository.GetTotalSpending(ctx, db.GetTotalSpendingParams{
		Ownerid:   user.ID,
		Startdate: filter.StartDate,
		Enddate:   filter.EndDate,
	})

	if err != nil {
		return nil, err
	}

	result := &gen.SpendingStats{
		Total: utils.FormatCurrencyFloat64(int32(incomeTotal.(int64))),
	}

	totalCount, err := r.Repository.CountSpendingTransactions(ctx, db.CountSpendingTransactionsParams{
		Ownerid:   user.ID,
		Startdate: filter.StartDate,
		Enddate:   filter.EndDate,
	})

	if err != nil {
		return result, nil
	}

	paginator := paging.NewOffsetPaginator(pageArgs, totalCount)
	transactionsResult := &gen.TransactionConnection{
		PageInfo: &paginator.PageInfo,
	}
	start := int32(paginator.Offset)
	limit := calculatePageLimit(pageArgs)

	incomeTransactions, err := r.Repository.ListIncomeTransactions(ctx, db.ListIncomeTransactionsParams{
		Ownerid:   user.ID,
		Startdate: filter.StartDate,
		Enddate:   filter.EndDate,
		Limit:     limit,
		Start:     start,
	})

	if err != nil {
		return nil, err
	}

	for i, row := range incomeTransactions {
		transactionsResult.Edges = append(transactionsResult.Edges, gen.TransactionEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	result.Transactions = transactionsResult

	return result, err
}

func (r *queryResolver) Income(ctx context.Context, input gen.StatsInput) (*gen.IncomeStats, error) {
	user := auth.GetCurrentUser(ctx)
	pageArgs := getPageArgs(ctx, "transactions")
	filter, err := parseStatsFilter(input.Filter)

	if err != nil {
		return nil, err
	}

	incomeTotal, err := r.Repository.GetTotalIncome(ctx, db.GetTotalIncomeParams{
		Ownerid:   user.ID,
		Startdate: filter.StartDate,
		Enddate:   filter.EndDate,
	})

	if err != nil {
		return nil, err
	}

	result := &gen.IncomeStats{
		Total: utils.FormatCurrencyFloat64(int32(incomeTotal.(int64))),
	}

	totalCount, err := r.Repository.CountIncomeTransactions(ctx, db.CountIncomeTransactionsParams{
		Ownerid:   user.ID,
		Startdate: filter.StartDate,
		Enddate:   filter.EndDate,
	})

	if err != nil {
		return result, nil
	}

	paginator := paging.NewOffsetPaginator(pageArgs, totalCount)
	transactionsResult := &gen.TransactionConnection{
		PageInfo: &paginator.PageInfo,
	}
	start := int32(paginator.Offset)
	limit := calculatePageLimit(pageArgs)

	incomeTransactions, err := r.Repository.ListIncomeTransactions(ctx, db.ListIncomeTransactionsParams{
		Ownerid:   user.ID,
		Startdate: filter.StartDate,
		Enddate:   filter.EndDate,
		Limit:     limit,
		Start:     start,
	})

	if err != nil {
		return nil, err
	}

	for i, row := range incomeTransactions {
		transactionsResult.Edges = append(transactionsResult.Edges, gen.TransactionEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	result.Transactions = transactionsResult

	return result, err
}

func (r *queryResolver) Net(ctx context.Context, input gen.StatsInput) (*gen.NetStats, error) {
	user := auth.GetCurrentUser(ctx)
	pageArgs := getPageArgs(ctx, "transactions")
	filter, err := parseStatsFilter(input.Filter)

	if err != nil {
		return nil, err
	}

	incomeTotal, err := r.Repository.GetNetIncome(ctx, db.GetNetIncomeParams{
		Ownerid:   user.ID,
		Startdate: filter.StartDate,
		Enddate:   filter.EndDate,
	})

	if err != nil {
		return nil, err
	}

	result := &gen.NetStats{
		Total: utils.FormatCurrencyFloat64(int32(incomeTotal.(int64))),
	}

	totalCount, err := r.Repository.CountTransactionsByDates(ctx, db.CountTransactionsByDatesParams{
		Ownerid:   user.ID,
		Startdate: filter.StartDate,
		Enddate:   filter.EndDate,
	})

	if err != nil {
		return result, nil
	}

	paginator := paging.NewOffsetPaginator(pageArgs, totalCount)
	transactionsResult := &gen.TransactionConnection{
		PageInfo: &paginator.PageInfo,
	}
	start := int32(paginator.Offset)
	limit := calculatePageLimit(pageArgs)

	incomeTransactions, err := r.Repository.ListTransactionsByDates(ctx, db.ListTransactionsByDatesParams{
		Ownerid:   user.ID,
		Startdate: filter.StartDate,
		Enddate:   filter.EndDate,
		Limit:     limit,
		Start:     start,
	})

	if err != nil {
		return nil, err
	}

	for i, row := range incomeTransactions {
		transactionsResult.Edges = append(transactionsResult.Edges, gen.TransactionEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	result.Transactions = transactionsResult

	return result, err
}
