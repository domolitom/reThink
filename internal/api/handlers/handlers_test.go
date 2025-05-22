// handlers_test.go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/domolitom/reThink/internal/models"
	"github.com/domolitom/reThink/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock DB for testing
type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) CreatePrediction(prediction *models.Prediction) error {
	args := m.Called(prediction)
	return args.Error(0)
}

func (m *MockDB) GetPredictions(userID int, page int, limit int) ([]models.Prediction, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]models.Prediction), args.Error(1)
}

func (m *MockDB) VotePrediction(vote *models.Vote) error {
	args := m.Called(vote)
	return args.Error(0)
}

// SetupTestRouter creates a Gin router for testing
func SetupTestRouter(handler *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Auth routes
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)

	// Prediction routes
	auth := r.Group("/")
	auth.Use(handler.AuthMiddleware())
	{
		auth.POST("/predictions", handler.CreatePrediction)
		auth.GET("/predictions", handler.GetPredictions)
		auth.POST("/predictions/:id/vote", handler.VotePrediction)
	}

	return r
}

func TestRegisterHandler(t *testing.T) {
	mockDB := new(MockDB)
	handler := NewHandler(mockDB)
	router := SetupTestRouter(handler)

	tests := []struct {
		name           string
		payload        models.RegisterRequest
		setupMock      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful Registration",
			payload: models.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockDB.On("GetUserByEmail", "test@example.com").Return(nil, utils.ErrNotFound)
				mockDB.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "User registered successfully",
			},
		},
		{
			name: "User Already Exists",
			payload: models.RegisterRequest{
				Name:     "Existing User",
				Email:    "existing@example.com",
				Password: "password123",
			},
			setupMock: func() {
				existingUser := &models.User{
					ID:       1,
					Name:     "Existing User",
					Email:    "existing@example.com",
					Password: "hashedpassword",
				}
				mockDB.On("GetUserByEmail", "existing@example.com").Return(existingUser, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"error": "User already exists",
			},
		},
		{
			name: "Invalid Email",
			payload: models.RegisterRequest{
				Name:     "Invalid Email User",
				Email:    "notavalidemail",
				Password: "password123",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid email format",
			},
		},
		{
			name: "Password Too Short",
			payload: models.RegisterRequest{
				Name:     "Short Password User",
				Email:    "valid@example.com",
				Password: "short",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Password must be at least 8 characters long",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock()

			// Create request
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			// Clear mock expectations
			mockDB.AssertExpectations(t)
			mockDB.ExpectedCalls = nil
		})
	}
}

func TestLoginHandler(t *testing.T) {
	mockDB := new(MockDB)
	handler := NewHandler(mockDB)
	router := SetupTestRouter(handler)

	// Hash a password for our test
	hashedPassword, _ := utils.HashPassword("password123")

	tests := []struct {
		name           string
		payload        models.LoginRequest
		setupMock      func()
		expectedStatus int
		checkToken     bool
	}{
		{
			name: "Successful Login",
			payload: models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				user := &models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: hashedPassword,
				}
				mockDB.On("GetUserByEmail", "test@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			checkToken:     true,
		},
		{
			name: "User Not Found",
			payload: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockDB.On("GetUserByEmail", "nonexistent@example.com").Return(nil, utils.ErrNotFound)
			},
			expectedStatus: http.StatusUnauthorized,
			checkToken:     false,
		},
		{
			name: "Incorrect Password",
			payload: models.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func() {
				user := &models.User{
					ID:       1,
					Email:    "test@example.com",
					Password: hashedPassword,
				}
				mockDB.On("GetUserByEmail", "test@example.com").Return(user, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			checkToken:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock()

			// Create request
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkToken {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				// Verify token exists
				assert.Contains(t, response, "token")
				assert.NotEmpty(t, response["token"])
			}

			// Clear mock expectations
			mockDB.AssertExpectations(t)
			mockDB.ExpectedCalls = nil
		})
	}
}

func TestCreatePredictionHandler(t *testing.T) {
	mockDB := new(MockDB)
	handler := NewHandler(mockDB)
	router := SetupTestRouter(handler)

	// Create a valid token for testing
	token, _ := utils.GenerateToken(1)

	tests := []struct {
		name           string
		payload        models.PredictionRequest
		setupMock      func()
		expectedStatus int
		token          string
	}{
		{
			name: "Successful Prediction Creation",
			payload: models.PredictionRequest{
				Title:       "Stock Market Prediction",
				Description: "I predict the S&P 500 will reach 5000 by end of year",
				Category:    "Finance",
				EndDate:     time.Now().AddDate(0, 3, 0),
			},
			setupMock: func() {
				mockDB.On("CreatePrediction", mock.AnythingOfType("*models.Prediction")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			token:          token,
		},
		{
			name: "Missing Title",
			payload: models.PredictionRequest{
				Title:       "",
				Description: "Some description",
				Category:    "Finance",
				EndDate:     time.Now().AddDate(0, 3, 0),
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			token:          token,
		},
		{
			name: "Invalid End Date (Past)",
			payload: models.PredictionRequest{
				Title:       "Past Prediction",
				Description: "This prediction has an end date in the past",
				Category:    "Technology",
				EndDate:     time.Now().AddDate(0, 0, -1),
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			token:          token,
		},
		{
			name: "Unauthorized Request",
			payload: models.PredictionRequest{
				Title:       "Valid Prediction",
				Description: "Some description",
				Category:    "Finance",
				EndDate:     time.Now().AddDate(0, 1, 0),
			},
			setupMock:      func() {},
			expectedStatus: http.StatusUnauthorized,
			token:          "invalid-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock()

			// Create request
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/predictions", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Clear mock expectations
			mockDB.AssertExpectations(t)
			mockDB.ExpectedCalls = nil
		})
	}
}

func TestGetPredictionsHandler(t *testing.T) {
	mockDB := new(MockDB)
	handler := NewHandler(mockDB)
	router := SetupTestRouter(handler)

	// Create a valid token for testing
	token, _ := utils.GenerateToken(1)

	// Sample predictions for response
	predictions := []models.Prediction{
		{
			ID:          1,
			UserID:      1,
			Title:       "Prediction 1",
			Description: "Description 1",
			Category:    "Finance",
			CreatedAt:   time.Now(),
			EndDate:     time.Now().AddDate(0, 3, 0),
		},
		{
			ID:          2,
			UserID:      1,
			Title:       "Prediction 2",
			Description: "Description 2",
			Category:    "Technology",
			CreatedAt:   time.Now(),
			EndDate:     time.Now().AddDate(0, 6, 0),
		},
	}

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedCount  int
		token          string
		queryParams    string
	}{
		{
			name: "Get Predictions - Default Page and Limit",
			setupMock: func() {
				mockDB.On("GetPredictions", 1, 1, 10).Return(predictions, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			token:          token,
			queryParams:    "",
		},
		{
			name: "Get Predictions - Custom Page and Limit",
			setupMock: func() {
				mockDB.On("GetPredictions", 1, 2, 5).Return(predictions[:1], nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			token:          token,
			queryParams:    "?page=2&limit=5",
		},
		{
			name:           "Unauthorized Request",
			setupMock:      func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedCount:  0,
			token:          "invalid-token",
			queryParams:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock()

			// Create request
			req, _ := http.NewRequest("GET", "/predictions"+tt.queryParams, nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response struct {
					Predictions []models.Prediction `json:"predictions"`
				}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Len(t, response.Predictions, tt.expectedCount)
			}

			// Clear mock expectations
			mockDB.AssertExpectations(t)
			mockDB.ExpectedCalls = nil
		})
	}
}

func TestVotePredictionHandler(t *testing.T) {
	mockDB := new(MockDB)
	handler := NewHandler(mockDB)
	router := SetupTestRouter(handler)

	// Create a valid token for testing
	token, _ := utils.GenerateToken(1)

	tests := []struct {
		name           string
		predictionID   string
		payload        models.VoteRequest
		setupMock      func()
		expectedStatus int
		token          string
	}{
		{
			name:         "Successful Vote",
			predictionID: "1",
			payload: models.VoteRequest{
				Value: true, // Agree with prediction
			},
			setupMock: func() {
				mockDB.On("VotePrediction", mock.AnythingOfType("*models.Vote")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			token:          token,
		},
		{
			name:         "Invalid Prediction ID",
			predictionID: "abc", // Not a number
			payload: models.VoteRequest{
				Value: false,
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			token:          token,
		},
		{
			name:         "Unauthorized Request",
			predictionID: "1",
			payload: models.VoteRequest{
				Value: true,
			},
			setupMock:      func() {},
			expectedStatus: http.StatusUnauthorized,
			token:          "invalid-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock()

			// Create request
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/predictions/"+tt.predictionID+"/vote", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tt.token)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Clear mock expectations
			mockDB.AssertExpectations(t)
			mockDB.ExpectedCalls = nil
		})
	}
}
