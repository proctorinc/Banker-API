package resolvers

import (
	"context"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	"github.com/proctorinc/banker/internal/graphql/utils"
)

type accountResolver struct{ *Resolver }

func (r *accountResolver) ID(ctx context.Context, account *db.Account) (uuid.UUID, error) {
	return account.ID, nil
}

func (r *accountResolver) SourceId(ctx context.Context, account *db.Account) (string, error) {
	masked := utils.MaskData(account.Sourceid)
	return masked, nil
}

func (r *accountResolver) UploadSource(ctx context.Context, account *db.Account) (string, error) {
	return string(account.Uploadsource), nil
}

func (r *accountResolver) Type(ctx context.Context, account *db.Account) (string, error) {
	return string(account.Type), nil
}

func (r *accountResolver) Name(ctx context.Context, account *db.Account) (string, error) {
	return account.Name, nil
}

func (r *accountResolver) RoutingNumber(ctx context.Context, account *db.Account) (*string, error) {
	if len(account.Routingnumber.String) > 0 {
		masked := utils.MaskData(account.Routingnumber.String)
		return &masked, nil
	}

	return nil, nil
}

// Queries

func (r *queryResolver) Account(ctx context.Context, accountId uuid.UUID) (*db.Account, error) {
	user := auth.GetCurrentUser(ctx)
	account, err := r.Repository.GetAccount(ctx, db.GetAccountParams{
		ID:      accountId,
		Ownerid: user.ID,
	})

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *queryResolver) Accounts(ctx context.Context) ([]db.Account, error) {
	user := auth.GetCurrentUser(ctx)

	return r.Repository.ListAccounts(ctx, user.ID)
}
