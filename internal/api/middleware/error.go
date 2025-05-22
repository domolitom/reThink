package middleware

import (
	"net/http"

	"github.com/domolitom/reThink/utils"
	"github.com/gin-gonic/gin"
)

// ErrorMiddleware handles and logs errors
func ErrorMiddleware(logger *utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("Request error:", err.Err.Error())
			}

			// If no response has been sent yet, send a generic error
			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "An error occurred while processing your request",
				})
			}
		}
	}
}
