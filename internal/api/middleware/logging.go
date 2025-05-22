package middleware

import (
	"time"

	"strconv"

	"github.com/domolitom/reThink/utils"
	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger *utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		startTime := time.Now()

		// Process request
		c.Next()

		// End time
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Log request details
		logger.Info(
			"Request:",
			"method="+c.Request.Method,
			"status="+strconv.Itoa(c.Writer.Status()),
			"status="+strconv.Itoa(c.Writer.Status()),
			"latency="+latency.String(),
			"ip="+c.ClientIP(),
		)
	}
}
