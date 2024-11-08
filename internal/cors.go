package cors

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{
		"http://localhost:5173",
	}

	corsConfig.AllowCredentials = true

	return cors.New(corsConfig)
}
