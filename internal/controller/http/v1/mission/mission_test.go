package mission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"DevelopsToday/internal/controller/http/middleware"
	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo/mocks"
	"DevelopsToday/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test constants
const (
	// Test mission data
	TestMissionID    = uint(1)
	InvalidMissionID = uint(999)
	TestCatID        = uint(4)
	InvalidCatID     = uint(999)

	// Test URLs
	MissionsBaseURL     = "/v1/missions"
	InvalidMissionIDStr = "invalid"

	// Test JSON
	InvalidJSON = "invalid json"
)

// MockLogger для тестування
type MockLogger struct{}

func (m *MockLogger) Debug(message interface{}, args ...interface{}) {}
func (m *MockLogger) Info(message string, args ...interface{})       {}
func (m *MockLogger) Warn(message string, args ...interface{})       {}
func (m *MockLogger) Error(message interface{}, args ...interface{}) {}
func (m *MockLogger) Fatal(message interface{}, args ...interface{}) {}

// Helper functions for tests
func createTestTargets(count int) []models.Target {
	targets := make([]models.Target, count)
	for i := 0; i < count; i++ {
		targets[i] = models.Target{
			Name:     fmt.Sprintf("Target %d", i+1),
			Country:  fmt.Sprintf("Country %d", i+1),
			Notes:    fmt.Sprintf("Notes %d", i+1),
			Complete: false,
		}
	}
	return targets
}

func makeJSONRequest(t *testing.T, method, url string, data interface{}) *http.Request {
	var jsonData []byte
	var err error

	if data != nil {
		jsonData, err = json.Marshal(data)
		assert.NoError(t, err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

func makeRequest(t *testing.T, method, url string) *http.Request {
	req, err := http.NewRequest(method, url, http.NoBody)
	assert.NoError(t, err)
	return req
}

func assertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, target interface{}) {
	assert.Equal(t, expectedStatus, w.Code)
	if target != nil {
		err := json.Unmarshal(w.Body.Bytes(), target)
		assert.NoError(t, err)
	}
}

func setupTestRouter() (*gin.Engine, *Service) {
	gin.SetMode(gin.TestMode)

	store := mocks.NewRepository()
	missionService := services.NewMission(store.Mission())
	mockLogger := &MockLogger{}

	service := NewImplService(missionService)
	handler := NewHandler(service, mockLogger)

	router := gin.New()
	router.Use(middleware.GlobalErrorHandler())
	v1 := router.Group("/v1")
	missions := v1.Group("/missions")
	{
		missions.POST("", handler.Create)
		missions.GET("", handler.List)
		missions.GET("/:id", handler.GetByID)
		missions.POST("/:id/assign", handler.AssignCat)
		missions.POST("/:id/complete", handler.MarkComplete)
		missions.DELETE("/:id", handler.Delete)
	}

	return router, service
}

func TestMissionController_Create(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should create mission with valid targets", func(t *testing.T) {
		createReq := CreateRequest{
			Targets: createTestTargets(2),
		}

		req := makeJSONRequest(t, "POST", MissionsBaseURL, createReq)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response models.Mission
		assertJSONResponse(t, w, http.StatusCreated, &response)
		assert.NotZero(t, response.ID)
		assert.Equal(t, 2, len(response.Targets))
	})

	t.Run("should fail with no targets", func(t *testing.T) {
		createReq := CreateRequest{
			Targets: []models.Target{},
		}

		req := makeJSONRequest(t, "POST", MissionsBaseURL, createReq)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail with more than 3 targets", func(t *testing.T) {
		createReq := CreateRequest{
			Targets: createTestTargets(4), // More than 3 targets
		}

		req := makeJSONRequest(t, "POST", MissionsBaseURL, createReq)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		req, err := http.NewRequest("POST", MissionsBaseURL, bytes.NewBuffer([]byte(InvalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})
}

func TestMissionController_List(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should return list of missions", func(t *testing.T) {
		req := makeRequest(t, "GET", MissionsBaseURL)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response []models.Mission
		assertJSONResponse(t, w, http.StatusOK, &response)
		assert.GreaterOrEqual(t, len(response), 5) // Should have at least 5 initial missions
	})
}

func TestMissionController_GetByID(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should return mission by valid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/missions/1", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Mission
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.GreaterOrEqual(t, len(response.Targets), 1)
	})

	t.Run("should return 404 for non-existing mission", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/missions/999", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestMissionController_AssignCat(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should assign cat to mission", func(t *testing.T) {
		assignReq := AssignCatRequest{CatID: 4} // Assign Felix
		jsonData, _ := json.Marshal(assignReq)

		req, _ := http.NewRequest("POST", "/v1/missions/4/assign", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/missions/1/assign", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail for non-existing cat", func(t *testing.T) {
		assignReq := AssignCatRequest{CatID: 999}
		jsonData, _ := json.Marshal(assignReq)

		req, _ := http.NewRequest("POST", "/v1/missions/4/assign", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestMissionController_MarkComplete(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should fail for mission with incomplete targets", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/missions/1/complete", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should succeed for mission with all targets complete", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/missions/2/complete", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should fail for non-existing mission", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/missions/999/complete", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/missions/invalid/complete", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMissionController_Delete(t *testing.T) {
	router, service := setupTestRouter()

	t.Run("should fail for assigned mission", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/missions/1", http.NoBody) // Mission 1 is assigned

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should succeed for unassigned mission", func(t *testing.T) {
		// First create an unassigned mission
		mission := &models.Mission{
			Complete: false,
			CatID:    nil,
			Targets: []models.Target{
				{Name: "To Delete", Country: "Test", Notes: "Test", Complete: false},
			},
		}
		err := service._missionContext.Create(context.TODO(), mission)
		assert.NoError(t, err)

		req, _ := http.NewRequest("DELETE", "/v1/missions/"+strconv.Itoa(int(mission.ID)), http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("should return 404 for non-existing mission", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/missions/999", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 400 for invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/missions/invalid", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
