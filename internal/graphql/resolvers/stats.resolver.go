package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/utils"
)

func (r *queryResolver) Stats(ctx context.Context, filter gen.StatsInput) (*gen.StatsResponse, error) {
	user := auth.GetCurrentUser(ctx)
	startDate, err := time.Parse(time.RFC3339, filter.StartDate)

	if err != nil {
		return nil, fmt.Errorf("Invalid date format. RFC3339 required")
	}

	endDate, err := time.Parse(time.RFC3339, filter.EndDate)

	if err != nil {
		return nil, fmt.Errorf("Invalid date format. RFC3339 required")
	}

	incomeTotal, err := r.Repository.GetTotalIncome(ctx, db.GetTotalIncomeParams{
		Ownerid:   user.ID,
		Startdate: startDate,
		Enddate:   endDate,
	})

	if err != nil {
		return nil, err
	}

	incomeTransactions, err := r.Repository.ListIncomeTransactions(ctx, db.ListIncomeTransactionsParams{
		Ownerid:   user.ID,
		Startdate: startDate,
		Enddate:   endDate,
	})

	if err != nil {
		return nil, err
	}

	income := gen.IncomeStats{
		Total:        utils.FormatCurrencyFloat64(int32(incomeTotal.(int64))),
		Transactions: incomeTransactions,
	}

	spendingTotal, err := r.Repository.GetTotalSpending(ctx, db.GetTotalSpendingParams{
		Ownerid:   user.ID,
		Startdate: startDate,
		Enddate:   endDate,
	})

	if err != nil {
		return nil, err
	}

	spendingTransactions, err := r.Repository.ListSpendingTransactions(ctx, db.ListSpendingTransactionsParams{
		Ownerid:   user.ID,
		Startdate: startDate,
		Enddate:   endDate,
	})

	if err != nil {
		return nil, err
	}

	spending := gen.SpendingStats{
		Total:        utils.FormatCurrencyFloat64(int32(spendingTotal.(int64))),
		Transactions: spendingTransactions,
	}

	net := gen.NetStats{
		Total: utils.FormatCurrencyFloat64(int32(incomeTotal.(int64)) + int32(spendingTotal.(int64))),
	}

	response := &gen.StatsResponse{
		Spending: &spending,
		Income:   &income,
		Net:      &net,
	}

	return response, nil
}
