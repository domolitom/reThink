package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/domolitom/reThink/internal/database"
	"github.com/domolitom/reThink/internal/models"
	"github.com/gin-gonic/gin"
)

type CreateMarketInput struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CloseDate   time.Time `json:"close_date" binding:"required"`
	ResolveDate time.Time `json:"resolve_date" binding:"required"`
}

type UpdateMarketInput struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CloseDate   time.Time `json:"close_date"`
	ResolveDate time.Time `json:"resolve_date"`
}

type ResolveMarketInput struct {
	Outcome bool `json:"outcome" binding:"required"`
}

// GetMarkets returns all markets with pagination
func GetMarkets(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Get status filter if provided
	status := c.Query("status")

	var markets []models.Market
	var total int64
	query := database.DB.Model(&models.Market{})

	// Apply status filter if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total records for pagination
	query.Count(&total)

	// Execute query with pagination
	result := query.Preload("Creator").Order("created_at desc").Limit(limit).Offset(offset).Find(&markets)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve markets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"markets": markets,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
			"pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetMarket returns a specific market by ID
func GetMarket(c *gin.Context) {
	id := c.Param("id")

	var market models.Market
	result := database.DB.Preload("Creator").First(&market, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Market not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"market": market})
}

// CreateMarket creates a new market
func CreateMarket(c *gin.Context) {
	var input CreateMarketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the current user ID from context
	userID, _ := c.Get("userID")

	// Create new market
	market := models.Market{
		Title:       input.Title,
		Description: input.Description,
		CreatorID:   userID.(uint),
		CloseDate:   input.CloseDate,
		ResolveDate: input.ResolveDate,
		Status:      models.MarketOpen,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if result := database.DB.Create(&market); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create market"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Market created successfully",
		"market":  market,
	})
}

// UpdateMarket updates an existing market
func UpdateMarket(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var market models.Market
	if result := database.DB.First(&market, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Market not found"})
		return
	}

	// Check if user is the creator
	if market.CreatorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update markets you created"})
		return
	}

	// Check if market is already resolved
	if market.Status == models.MarketResolved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update a resolved market"})
		return
	}

	var input UpdateMarketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if input.Title != "" {
		market.Title = input.Title
	}
	if input.Description != "" {
		market.Description = input.Description
	}
	if !input.CloseDate.IsZero() {
		market.CloseDate = input.CloseDate
	}
	if !input.ResolveDate.IsZero() {
		market.ResolveDate = input.ResolveDate
	}

	market.UpdatedAt = time.Now()

	if result := database.DB.Save(&market); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update market"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Market updated successfully",
		"market":  market,
	})
}

// ResolveMarket resolves a market with a final outcome
func ResolveMarket(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")

	var market models.Market
	if result := database.DB.First(&market, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Market not found"})
		return
	}

	// Check if user is the creator
	if market.CreatorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the creator can resolve this market"})
		return
	}

	// Check if market is already resolved
	if market.Status == models.MarketResolved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Market is already resolved"})
		return
	}

	var input ResolveMarketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update market status and outcome
	outcome := input.Outcome
	market.Status = models.MarketResolved
	market.Outcome = &outcome
	market.UpdatedAt = time.Now()

	// Start a transaction
	tx := database.DB.Begin()

	// Save the market resolution
	if result := tx.Save(&market); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve market"})
		return
	}

	// Update user prediction scores
	var predictions []models.Prediction
	if result := tx.Where("market_id = ?", market.ID).Preload("User").Find(&predictions); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve predictions"})
		return
	}

	// Calculate and update scores for each user who predicted
	for _, prediction := range predictions {
		// Calculate score adjustment based on prediction accuracy
		// This is a basic scoring algorithm; you might want more sophisticated ones
		var scoreAdjustment float64
		if prediction.Prediction == outcome {
			// Correct prediction: higher confidence = higher score
			scoreAdjustment = prediction.Confidence / 50.0 // Normalized to 0-2
		} else {
			// Wrong prediction: higher confidence = larger penalty
			scoreAdjustment = -prediction.Confidence / 100.0 // Normalized to 0-1
		}

		// Update user's prediction score
		prediction.User.PredictionScore += scoreAdjustment
		if result := tx.Save(&prediction.User); result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user scores"})
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete market resolution"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Market resolved successfully",
		"market":  market,
	})
}
