// middleware_test.go
package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/domolitom/reThink/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin engine
	r := gin.New()

	// Create auth middleware
	authMiddleware := AuthMiddleware()

	// Create a test endpoint with the middleware
	r.GET("/protected", authMiddleware, func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if exists {
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		}
	})

	// Generate a valid token for testing
	userID := 123
	token, err := utils.GenerateToken(userID)
	assert.NoError(t, err)

	// Test cases
	tests := []struct {
		name           string
		token          string
		expectedStatus int
		checkUserID    bool
	}{
		{
			name:           "Valid Token",
			token:          token,
			expectedStatus: http.StatusOK,
			checkUserID:    true,
		},
		{
			name:           "No Token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			checkUserID:    false,
		},
		{
			name:           "Invalid Token Format",
			token:          "invalid.token.format",
			expectedStatus: http.StatusUnauthorized,
			checkUserID:    false,
		},
		{
			name:           "Token Without Bearer Prefix",
			token:          token, // But we'll send it without "Bearer "
			expectedStatus: http.StatusUnauthorized,
			checkUserID:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/protected", nil)

			// Add Authorization header if token is provided
			if tt.token != "" {
				if tt.name == "Token Without Bearer Prefix" {
					req.Header.Set("Authorization", tt.token)
				} else {
					req.Header.Set("Authorization", "Bearer "+tt.token)
				}
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkUserID {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, float64(userID), response["user_id"])
			}
		})
	}
}
