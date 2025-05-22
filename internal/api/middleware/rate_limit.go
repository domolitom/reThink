package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Simple in-memory rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
	}
}

// RateLimitMiddleware limits the number of requests per IP in a given time window
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		// Clean up old requests
		now := time.Now()
		cutoff := now.Add(-window)

		var recent []time.Time
		for _, t := range limiter.requests[ip] {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}

		// Update requests with only recent ones
		limiter.requests[ip] = recent

		// Check if request limit is exceeded
		if len(recent) >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		// Add current request
		limiter.requests[ip] = append(limiter.requests[ip], now)

		c.Next()
	}
}
