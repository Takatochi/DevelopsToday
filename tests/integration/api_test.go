//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"DevelopsToday/config"
	"DevelopsToday/internal/app"
	"DevelopsToday/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) (*gin.Engine, func()) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test configuration
	cfg := &config.Config{
		HTTP: config.HTTP{
			Port: "8080",
		},
		PG: config.PG{
			URL:     "postgres://spy_cats:secret_password@127.0.0.1:5432/spy_cats?sslmode=disable",
			PoolMax: 10,
		},
		JWT: config.JWT{
			Secret:          "test-secret-key-for-integration-tests",
			AccessTokenTTL:  900,    // 15 minutes
			RefreshTokenTTL: 604800, // 7 days
		},
		App: config.App{
			Name:    "spy-cats-api",
			Version: "1.0.0",
		},
		Log: config.Log{
			Level: "info",
		},
		Cache: config.Cache{
			Type: "memory",
		},
		Swagger: config.Swagger{
			Enabled: false,
		},
	}

	// Initialize app
	application, err := app.New(cfg)
	require.NoError(t, err)

	cleanup := func() {
		// Clean up test data if needed
	}

	return application.Handler, cleanup
}

func TestAuthIntegration(t *testing.T) {
	router, cleanup := setupTestApp(t)
	defer cleanup()

	t.Run("Login with valid credentials should return tokens", func(t *testing.T) {
		loginData := map[string]string{
			"username": "admin",
			"password": "admin123",
		}

		body, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Contains(t, response, "access_token")
		assert.Contains(t, response, "refresh_token")
		assert.Contains(t, response, "user")
	})

	t.Run("Login with invalid credentials should return error", func(t *testing.T) {
		loginData := map[string]string{
			"username": "admin",
			"password": "wrongpassword",
		}

		body, _ := json.Marshal(loginData)
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCatsIntegration(t *testing.T) {
	router, cleanup := setupTestApp(t)
	defer cleanup()

	// First, get auth token
	token := getAuthToken(t, router)

	t.Run("GET /v1/cats should return cats list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/cats", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var cats []models.Cat
		err := json.Unmarshal(w.Body.Bytes(), &cats)
		require.NoError(t, err)
		assert.IsType(t, []models.Cat{}, cats)
	})

	t.Run("POST /v1/cats should create new cat", func(t *testing.T) {
		newCat := map[string]interface{}{
			"name":       "TestCat",
			"breed":      "Persian",
			"experience": 3,
			"salary":     1000,
		}

		body, _ := json.Marshal(newCat)
		req := httptest.NewRequest(http.MethodPost, "/v1/cats", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createdCat models.Cat
		err := json.Unmarshal(w.Body.Bytes(), &createdCat)
		require.NoError(t, err)
		assert.Equal(t, "TestCat", createdCat.Name)
		assert.Equal(t, "Persian", createdCat.Breed)
	})

	t.Run("GET /v1/cats without auth should return 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/cats", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestMissionsIntegration(t *testing.T) {
	router, cleanup := setupTestApp(t)
	defer cleanup()

	token := getAuthToken(t, router)

	t.Run("GET /v1/missions should return missions list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/missions", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var missions []models.Mission
		err := json.Unmarshal(w.Body.Bytes(), &missions)
		require.NoError(t, err)
		assert.IsType(t, []models.Mission{}, missions)
	})

	t.Run("POST /v1/missions should create new mission", func(t *testing.T) {
		newMission := map[string]interface{}{
			"targets": []map[string]interface{}{
				{
					"name":    "Test Target",
					"country": "Test Country",
					"notes":   "Test notes",
				},
			},
		}

		body, _ := json.Marshal(newMission)
		req := httptest.NewRequest(http.MethodPost, "/v1/missions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createdMission models.Mission
		err := json.Unmarshal(w.Body.Bytes(), &createdMission)
		require.NoError(t, err)
		assert.Len(t, createdMission.Targets, 1)
		assert.Equal(t, "Test Target", createdMission.Targets[0].Name)
	})
}

func TestHealthCheck(t *testing.T) {
	router, cleanup := setupTestApp(t)
	defer cleanup()

	t.Run("GET /health should return OK", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "OK", response["status"])
	})
}

// Helper function to get auth token for tests
func getAuthToken(t *testing.T, router *gin.Engine) string {
	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}

	body, _ := json.Marshal(loginData)
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	token, ok := response["access_token"].(string)
	require.True(t, ok)
	require.NotEmpty(t, token)

	return token
}
