package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/paging"
)

func (r *userResolver) ID(ctx context.Context, user *db.User) (string, error) {
	return user.ID.String(), nil
}

func (r *userResolver) Role(ctx context.Context, user *db.User) (string, error) {
	return string(user.Role), nil
}

func (r *userResolver) Accounts(ctx context.Context, user *db.User, page *paging.PageArgs) (*gen.AccountConnection, error) {
	totalCount, err := r.Repository.CountAccounts(ctx, user.ID)

	if err != nil {
		return &gen.AccountConnection{
			PageInfo: paging.NewEmptyPageInfo(),
		}, err
	}

	paginator := paging.NewOffsetPaginator(page, totalCount)
	result := &gen.AccountConnection{
		PageInfo: &paginator.PageInfo,
	}
	start := int32(paginator.Offset)
	limit := calculatePageLimit(page)

	accounts, err := r.Repository.ListAccounts(ctx, db.ListAccountsParams{
		Ownerid: user.ID,
		Limit:   limit,
		Start:   start,
	})

	for i, row := range accounts {
		result.Edges = append(result.Edges, gen.AccountEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	return result, err
}

func (r *userResolver) Transactions(ctx context.Context, user *db.User, page *paging.PageArgs) (*gen.TransactionConnection, error) {
	totalCount, err := r.Repository.CountTransactions(ctx, user.ID)

	if err != nil {
		return &gen.TransactionConnection{
			PageInfo: paging.NewEmptyPageInfo(),
		}, err
	}

	paginator := paging.NewOffsetPaginator(page, totalCount)
	result := &gen.TransactionConnection{
		PageInfo: &paginator.PageInfo,
	}
	start := int32(paginator.Offset)
	limit := calculatePageLimit(page)

	transactions, err := r.Repository.ListTransactions(ctx, db.ListTransactionsParams{
		Ownerid: user.ID,
		Limit:   limit,
		Start:   start,
	})

	for i, row := range transactions {
		result.Edges = append(result.Edges, gen.TransactionEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	return result, err
}

func (r *userResolver) Merchants(ctx context.Context, user *db.User, page *paging.PageArgs) (*gen.MerchantConnection, error) {
	totalCount, err := r.Repository.CountMerchants(ctx, user.ID)

	if err != nil {
		return &gen.MerchantConnection{
			PageInfo: paging.NewEmptyPageInfo(),
		}, err
	}

	paginator := paging.NewOffsetPaginator(page, totalCount)
	result := &gen.MerchantConnection{
		PageInfo: &paginator.PageInfo,
	}
	start := int32(paginator.Offset)
	limit := calculatePageLimit(page)

	merchants, err := r.Repository.ListMerchants(ctx, db.ListMerchantsParams{
		Ownerid: user.ID,
		Limit:   limit,
		Start:   start,
	})

	for i, row := range merchants {
		result.Edges = append(result.Edges, gen.MerchantEdge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   &row,
		})
	}

	return result, err
}

// Queries

func (r *queryResolver) Me(ctx context.Context) (*db.User, error) {
	user := auth.GetCurrentUser(ctx)

	if user == nil {
		return nil, fmt.Errorf("no authenticated user")
	}

	return user, nil
}

func (r *queryResolver) User(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	user, err := r.Repository.GetUser(ctx, userId)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Mutations

func (r *mutationResolver) Login(ctx context.Context, data gen.LoginInput) (*db.User, error) {
	return r.AuthService.Login(ctx, auth.LoginInput{
		Email:    data.Email,
		Password: data.Password,
	})
}

func (r *mutationResolver) Logout(ctx context.Context) (string, error) {
	return r.AuthService.Logout(ctx)
}

func (r *mutationResolver) Register(ctx context.Context, data gen.RegisterInput) (*db.User, error) {
	return r.AuthService.Register(ctx, auth.RegisterInput{
		Email:    data.Email,
		Username: data.Username,
		Password: data.Password,
	})
}

func (r *mutationResolver) DeleteUser(ctx context.Context) (*db.User, error) {
	user := auth.GetCurrentUser(ctx)
	deleted, err := r.Repository.DeleteUser(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	return &deleted, nil
}
