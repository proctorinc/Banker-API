package resolvers

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/99designs/gqlgen/graphql"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/paging"
)

type StatsFilter struct {
	StartDate time.Time
	EndDate   time.Time
}

func calculatePageLimit(page *paging.PageArgs) int32 {
	var limit float64 = paging.MAX_PAGE_SIZE

	if page != nil && page.First != nil {
		limit = math.Min(limit, float64(*page.First))
	}

	return int32(limit)
}

func parseStatsFilter(input *gen.DateFilter) (*StatsFilter, error) {
	startDate, err := time.Parse(time.RFC3339, input.StartDate)

	if err != nil {
		return nil, fmt.Errorf("Invalid date format. RFC3339 required")
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)

	if err != nil {
		return nil, fmt.Errorf("Invalid date format. RFC3339 required")
	}

	filter := &StatsFilter{
		StartDate: startDate,
		EndDate:   endDate,
	}

	return filter, nil
}

func getPageArgs(ctx context.Context, fieldName string) *paging.PageArgs {
	pageArgs := &paging.PageArgs{
		First: nil,
		After: nil,
	}

	for _, field := range graphql.CollectFieldsCtx(ctx, nil) {
		if field.Name == fieldName {
			requestVariables := graphql.GetOperationContext(ctx).Variables
			args := field.ArgumentMap(requestVariables)

			if pageArgsMap, ok := args["page"].(map[string]any); ok {
				if first, ok := pageArgsMap["first"].(int64); ok {
					firstInt := int(first)
					pageArgs.First = &firstInt
				}

				if after, ok := pageArgsMap["after"].(string); ok {
					pageArgs.After = &after
				}
			}
		}
	}

	return pageArgs
}
