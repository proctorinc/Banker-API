package resolvers

import (
	"math"

	"github.com/proctorinc/banker/internal/graphql/paging"
)

func calculatePageLimit(page *paging.PageArgs) int32 {
	var limit float64 = paging.MAX_PAGE_SIZE

	if page != nil && page.First != nil {
		limit = math.Min(limit, float64(*page.First))
	}

	return int32(limit)
}
