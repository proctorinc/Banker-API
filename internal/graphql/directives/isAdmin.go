package directives

import (
	"context"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/proctorinc/banker/internal/auth"
)

func IsAdmin(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	user := auth.GetCurrentUser(ctx)

	if user != nil && user.Role == "ADMIN" {
		log.Printf("IsAdmin? Yes: %s", user.ID)
		return next(ctx)
	}

	log.Printf("IsAdmin? No")
	return nil, fmt.Errorf("You must be an admin to request this endpoint")
}
