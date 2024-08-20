package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/db"
)

func GraphqlHandler(repo db.Repository) gin.HandlerFunc {
	h := handler.NewDefaultServer(
		NewExecutableSchema(Config{
			Resolvers: &Resolver{
				Repository: repo,
			},
		}),
	)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func NewPlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
