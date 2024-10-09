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

func (r *merchantResolver) OwnerId(ctx context.Context, merchant *db.Merchant) (string, error) {
	return merchant.Ownerid.String(), nil
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
