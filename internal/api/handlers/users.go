package handlers

import (
	"net/http"
	"strconv"

	"github.com/domolitom/reThink/internal/database"
	"github.com/domolitom/reThink/internal/models"
	"github.com/gin-gonic/gin"
)

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Don't expose password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUser returns a specific user by ID
func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if result := database.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Don't expose password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateCurrentUser updates the current user's profile
func UpdateCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Only allow updating certain fields
	type UpdateUserInput struct {
		Bio      string `json:"bio"`
		Username string `json:"username"`
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if new username is available if it's being changed
	if input.Username != "" && input.Username != user.Username {
		var existingUser models.User
		if result := database.DB.Where("username = ?", input.Username).First(&existingUser); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
			return
		}
		user.Username = input.Username
	}

	// Update bio if provided
	if input.Bio != "" {
		user.Bio = input.Bio
	}

	if result := database.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Don't expose password hash
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

// GetUserStats returns prediction statistics for a user
func GetUserStats(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if result := database.DB.First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get total number of predictions
	var totalPredictions int64
	database.DB.Model(&models.Prediction{}).Where("user_id = ?", id).Count(&totalPredictions)

	// Get number of predictions on resolved markets
	var resolvedPredictions []models.Prediction
	database.DB.Joins("JOIN markets ON predictions.market_id = markets.id").
		Where("predictions.user_id = ? AND markets.status = ?", id, models.MarketResolved).
		Find(&resolvedPredictions)

	// Calculate correct predictions
	correctPredictions := 0
	for _, pred := range resolvedPredictions {
		var market models.Market
		database.DB.First(&market, pred.MarketID)

		if market.Outcome != nil && pred.Prediction == *market.Outcome {
			correctPredictions++
		}
	}

	// Calculate accuracy if there are resolved predictions
	var accuracy float64
	if len(resolvedPredictions) > 0 {
		accuracy = float64(correctPredictions) / float64(len(resolvedPredictions)) * 100
	}

	// Get recent predictions
	var recentPredictions []models.Prediction
	database.DB.Where("user_id = ?", id).
		Preload("Market").
		Order("created_at desc").
		Limit(5).
		Find(&recentPredictions)

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"total_predictions":    totalPredictions,
			"resolved_predictions": len(resolvedPredictions),
			"correct_predictions":  correctPredictions,
			"accuracy":             accuracy,
			"prediction_score":     user.PredictionScore,
			"recent_predictions":   recentPredictions,
		},
	})
}

// GetLeaderboard returns top users by prediction score
func GetLeaderboard(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var users []models.User
	var total int64

	// Count total users
	database.DB.Model(&models.User{}).Count(&total)

	// Get top users by prediction score
	result := database.DB.Model(&models.User{}).
		Select("id, username, bio, prediction_score, created_at, updated_at").
		Order("prediction_score desc").
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": users,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}
