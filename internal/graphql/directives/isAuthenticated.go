package directives

import (
	"context"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/proctorinc/banker/internal/auth"
)

func IsAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := auth.GetCurrentUser(ctx)
	log.Printf("User %s", user)
	if user == nil {
		log.Printf("No user exists")
		return nil, fmt.Errorf("You must be authenticated to request this endpoint")
	}

	log.Printf("User exists!")
	return next(ctx)
}
