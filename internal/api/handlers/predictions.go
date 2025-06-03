package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/domolitom/reThink/internal/database"
	"github.com/domolitom/reThink/internal/models"
	"github.com/gin-gonic/gin"
)

type CreatePredictionInput struct {
	Prediction bool    `json:"prediction" binding:"required"`
	Confidence float64 `json:"confidence" binding:"required,min=0,max=100"`
}

type UpdatePredictionInput struct {
	Prediction bool    `json:"prediction"`
	Confidence float64 `json:"confidence" binding:"min=0,max=100"`
}

// GetMarketPredictions returns all predictions for a specific market
func GetMarketPredictions(c *gin.Context) {
	marketID := c.Param("id")

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var predictions []models.Prediction
	var total int64

	// Count total records for pagination
	database.DB.Model(&models.Prediction{}).Where("market_id = ?", marketID).Count(&total)

	// Execute query with pagination
	result := database.DB.Where("market_id = ?", marketID).
		Preload("User").
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&predictions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve predictions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"predictions": predictions,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// CreatePrediction adds a new prediction for a market
func CreatePrediction(c *gin.Context) {
	marketID := c.Param("id")
	userID, _ := c.Get("userID")

	// Check if market exists
	var market models.Market
	if result := database.DB.First(&market, marketID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Market not found"})
		return
	}

	// Check if market is open for predictions
	if market.Status != models.MarketOpen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Market is not open for predictions"})
		return
	}

	// Check if user already made a prediction for this market
	var existingPrediction models.Prediction
	result := database.DB.Where("market_id = ? AND user_id = ?", marketID, userID).First(&existingPrediction)

	// Parse input
	var input CreatePredictionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If user already made a prediction, update it
	if result.RowsAffected > 0 {
		existingPrediction.Prediction = input.Prediction
		existingPrediction.Confidence = input.Confidence
		existingPrediction.UpdatedAt = time.Now()

		if result := database.DB.Save(&existingPrediction); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prediction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Prediction updated successfully",
			"prediction": existingPrediction,
		})
		return
	}

	// Create new prediction
	marketIDUint, _ := strconv.ParseUint(marketID, 10, 32)
	// Adjust field names and types to match models.Prediction struct definition
	prediction := models.Prediction{
		UserID:     uint(userID.(int)), // or userID.(uint) if userID is already uint
		MarketID:   uint(marketIDUint),
		IsTrue:     input.Prediction, // Replace with the correct boolean field name
		Confidence: input.Confidence, // Replace with the correct confidence field name
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(), // Replace with the correct updated/modified field name
	}

	if result := database.DB.Create(&prediction); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prediction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Prediction created successfully",
		"prediction": prediction,
	})
}

// UpdatePrediction updates an existing prediction
func UpdatePrediction(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var prediction models.Prediction
	if result := database.DB.First(&prediction, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prediction not found"})
		return
	}

	// Check if the prediction belongs to the user
	if prediction.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own predictions"})
		return
	}

	// Get the associated market
	var market models.Market
	if result := database.DB.First(&market, prediction.MarketID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Market not found"})
		return
	}

	// Check if market is still open
	if market.Status != models.MarketOpen {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Market is not open for predictions"})
		return
	}

	var input UpdatePredictionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	prediction.Prediction = input.Prediction
	prediction.Confidence = input.Confidence
	prediction.UpdatedAt = time.Now()

	if result := database.DB.Save(&prediction); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prediction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Prediction updated successfully",
		"prediction": prediction,
	})
}
