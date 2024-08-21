package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/db"
	"github.com/proctorinc/banker/internal/graphql/directives"
)

func GraphqlHandler(repo db.Repository) gin.HandlerFunc {
	config := Config{
		Resolvers: &Resolver{
			Repository: repo,
		},
	}

	config.Directives.IsAuthenticated = directives.IsAuthenticated

	handler := handler.NewDefaultServer(
		NewExecutableSchema(config),
	)

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func NewPlaygroundHandler() gin.HandlerFunc {
	handler := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
