package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/proctorinc/banker/internal/db"
)

func NewHandler(repo db.Repository) http.Handler {
	return handler.NewDefaultServer(
		NewExecutableSchema(Config{
			Resolvers: &Resolver{
				Repository: repo,
			},
		}),
	)
}

func NewPlaygroundHandler(endpoint string) http.Handler {
	return playground.Handler("GraphQL Playground", endpoint)
}
