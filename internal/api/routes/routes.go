package routes

import (
	"github.com/domolitom/reThink/internal/api/handlers"
	"github.com/domolitom/reThink/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	// API routes with authentication
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// User routes
		api.GET("/users/me", handlers.GetCurrentUser)
		api.GET("/users/:id", handlers.GetUser)
		api.PUT("/users/me", handlers.UpdateCurrentUser)

		// Market routes
		api.GET("/markets", handlers.GetMarkets)
		api.GET("/markets/:id", handlers.GetMarket)
		api.POST("/markets", handlers.CreateMarket)
		api.PUT("/markets/:id", handlers.UpdateMarket)
		api.POST("/markets/:id/resolve", handlers.ResolveMarket)

		// Prediction routes
		api.GET("/markets/:id/predictions", handlers.GetMarketPredictions)
		api.POST("/markets/:id/predict", handlers.CreatePrediction)
		api.PUT("/predictions/:id", handlers.UpdatePrediction)

		// Stats routes
		api.GET("/users/:id/stats", handlers.GetUserStats)
		api.GET("/leaderboard", handlers.GetLeaderboard)
	}
}
