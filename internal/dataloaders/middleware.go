package dataloaders

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/proctorinc/banker/internal/db"
)

// Middleware stores Loaders as a request-scoped context value.
func Middleware(repo db.Repository) func(gin.HandlerFunc) gin.HandlerFunc {
	return func(next gin.HandlerFunc) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			loaders := newLoaders(ctx, repo)
			newCtx := context.WithValue(ctx.Request.Context(), key, loaders)
			ctx.Request = ctx.Request.WithContext(newCtx)
			ctx.Next()
			next(ctx)
		}
	}
}
