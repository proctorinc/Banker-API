package paging

import "context"

type PageInfoResolver interface {
	HasPreviousPage(ctx context.Context, pageInfo *PageInfo) (bool, error)
	HasNextPage(ctx context.Context, pageInfo *PageInfo) (bool, error)
	TotalCount(ctx context.Context, pageInfo *PageInfo) (*int, error)
	StartCursor(ctx context.Context, pageInfo *PageInfo) (*string, error)
	EndCursor(ctx context.Context, pageInfo *PageInfo) (*string, error)
}
