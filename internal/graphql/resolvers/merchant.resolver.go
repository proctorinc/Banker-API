package resolvers

import (
	"context"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
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

func (r *queryResolver) Merchants(ctx context.Context) ([]db.Merchant, error) {
	user := auth.GetCurrentUser(ctx)

	return r.Repository.ListMerchants(ctx, user.ID)
}
