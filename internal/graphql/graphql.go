package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/auth"
	"github.com/proctorinc/banker/internal/dataloaders"
	"github.com/proctorinc/banker/internal/db"
	"github.com/proctorinc/banker/internal/graphql/directives"
	gen "github.com/proctorinc/banker/internal/graphql/generated"
	"github.com/proctorinc/banker/internal/graphql/resolvers"
)

func GraphqlHandler(repo db.Repository, loaders dataloaders.Retriever) gin.HandlerFunc {
	config := gen.Config{
		Resolvers: &resolvers.Resolver{
			Repository:  repo,
			AuthService: *auth.NewAuthService(repo),
			DataLoaders: loaders,
		},
	}

	// Add GraphQL directives
	config.Directives.IsAuthenticated = directives.IsAuthenticated
	config.Directives.IsAdmin = directives.IsAdmin

	handler := handler.NewDefaultServer(
		gen.NewExecutableSchema(config),
	)

	// Limit queries to 5 levels of complexity
	handler.Use(extension.FixedComplexityLimit(50))

	handler.AddTransport(transport.MultipartForm{
		MaxUploadSize: 5 * 1_000_000, // 5 MB max
		MaxMemory:     5 * 1_000_000, // 5 MB max
	})

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
