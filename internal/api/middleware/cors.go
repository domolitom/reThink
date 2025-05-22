package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware configures CORS settings for the API
func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()

	// Allow all origins in development - should be restricted in production
	config.AllowAllOrigins = true

	// Allow credentials
	config.AllowCredentials = true

	// Allow common methods
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

	// Allow common headers
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
	}

	return cors.New(config)
}
