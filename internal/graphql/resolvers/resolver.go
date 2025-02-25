package resolvers

import (
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/dataloaders"
	"github.com/proctorinc/banker/internal/db"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
)

type Resolver struct {
	Repository  db.Repository
	AuthService auth.AuthService
	DataLoaders dataloaders.Retriever
}

// Base resolvers
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type pageInfoResolver struct{ *Resolver }

// Models
type userResolver struct{ *Resolver }
type accountResolver struct{ *Resolver }
type accountSyncItemResolver struct{ *Resolver }
type transactionResolver struct{ *Resolver }
type merchantResolver struct{ *Resolver }
type fundResolver struct{ *Resolver }
type fundAllocationResolver struct{ *Resolver }
type statsResolver struct{ *Resolver }
type monthsResolver struct{ *Resolver }
type fundsResponseResolver struct{ *Resolver }

func (r *Resolver) Mutation() gen.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() gen.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) User() gen.UserResolver {
	return &userResolver{r}
}

func (r *Resolver) Account() gen.AccountResolver {
	return &accountResolver{r}
}

func (r *Resolver) AccountSyncItem() gen.AccountSyncItemResolver {
	return &accountSyncItemResolver{r}
}

func (r *Resolver) Transaction() gen.TransactionResolver {
	return &transactionResolver{r}
}

func (r *Resolver) Merchant() gen.MerchantResolver {
	return &merchantResolver{r}
}

func (r *Resolver) Fund() gen.FundResolver {
	return &fundResolver{r}
}

func (r *Resolver) FundAllocation() gen.FundAllocationResolver {
	return &fundAllocationResolver{r}
}

func (r *Resolver) PageInfo() gen.PageInfoResolver {
	return &pageInfoResolver{r}
}

func (r *Resolver) FundsResponse() gen.FundsResponseResolver {
	return &fundsResponseResolver{r}
}
