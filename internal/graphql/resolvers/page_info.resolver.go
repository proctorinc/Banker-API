package resolvers

import (
	"context"

	"github.com/proctorinc/banker/internal/graphql/paging"
)

func (r *pageInfoResolver) TotalCount(ctx context.Context, pageInfo *paging.PageInfo) (*int, error) {
	return pageInfo.TotalCount()
}

func (r *pageInfoResolver) HasPreviousPage(ctx context.Context, pageInfo *paging.PageInfo) (bool, error) {
	return pageInfo.HasPreviousPage()
}

func (r *pageInfoResolver) HasNextPage(ctx context.Context, pageInfo *paging.PageInfo) (bool, error) {
	return pageInfo.HasNextPage()
}

func (r *pageInfoResolver) StartCursor(ctx context.Context, pageInfo *paging.PageInfo) (*string, error) {
	return pageInfo.StartCursor()
}

func (r *pageInfoResolver) EndCursor(ctx context.Context, pageInfo *paging.PageInfo) (*string, error) {
	return pageInfo.EndCursor()
}
