package resolvers

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/paging"
)

func (r *merchantResolver) ID(ctx context.Context, merchant *db.Merchant) (uuid.UUID, error) {
	return merchant.ID, nil
}

func (r *merchantResolver) Name(ctx context.Context, merchant *db.Merchant) (string, error) {
	return merchant.Name, nil
}

func (r *merchantResolver) SourceID(ctx context.Context, merchant *db.Merchant) (*string, error) {
	if merchant.Sourceid.Valid {
		return &merchant.Sourceid.String, nil
	}

	return nil, nil
}

func (r *merchantResolver) OwnerId(ctx context.Context, merchant *db.Merchant) (string, error) {
	return merchant.Ownerid.String(), nil
}

func (r *merchantResolver) Transactions(ctx context.Context, merchant *db.Merchant) ([]db.Transaction, error) {
	return r.Repository.ListTransactionsByMerchantId(ctx, db.ListTransactionsByMerchantIdParams{
		Ownerid:    merchant.Ownerid,
		Merchantid: merchant.ID,
	})
}

// Queries

func (r *queryResolver) Merchant(ctx context.Context, merchantId uuid.UUID) (*db.Merchant, error) {
	user := auth.GetCurrentUser(ctx)
	merchant, err := r.Repository.GetMerchant(ctx, db.GetMerchantParams{
		ID:      merchantId,
		Ownerid: user.ID,
	})

	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (r *queryResolver) Merchants(ctx context.Context, page *paging.PageArgs) (*gen.MerchantConnection, error) {
	user := auth.GetCurrentUser(ctx)

	totalCount, err := r.Repository.ListMerchantsCount(ctx, user.ID)

	if err != nil {
		return &gen.MerchantConnection{
			PageInfo: paging.NewEmptyPageInfo(),
		}, err
	}

	var limit float64 = float64(totalCount)

	if page != nil && page.First != nil {
		limit = math.Min(float64(*page.First), limit)
	}

	paginator := paging.NewOffsetPaginator(page, totalCount)

	result := &gen.MerchantConnection{
		PageInfo: &paginator.PageInfo,
	}

	merchants, err := r.Repository.ListMerchants(ctx, db.ListMerchantsParams{
		Ownerid: user.ID,
		Limit:   int32(math.Min(limit, paging.MAX_PAGE_SIZE)),
		Start:   int32(paginator.Offset),
	})

	for i, row := range merchants {
		result.Edges = append(result.Edges, gen.MerchantEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	return result, err
}
