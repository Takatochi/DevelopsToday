package target

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
	// Test target data
	TestMissionID      = uint(1)
	TestTargetID       = uint(1)
	CompletedMissionID = uint(2)
	CompletedTargetID  = uint(3)
	InvalidMissionID   = uint(999)
	InvalidTargetID    = uint(999)

	// Test URLs
	TargetsBaseURL      = "/v1/missions"
	InvalidMissionIDStr = "invalid"
	InvalidTargetIDStr  = "invalid"

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
func createTestTarget(name, country, notes string) models.Target {
	return models.Target{
		Name:     name,
		Country:  country,
		Notes:    notes,
		Complete: false,
	}
}

func makePOSTRequest(t *testing.T, url string, data interface{}) *http.Request {
	var jsonData []byte
	var err error

	if data != nil {
		jsonData, err = json.Marshal(data)
		assert.NoError(t, err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

func makePUTRequest(t *testing.T, url string, data interface{}) *http.Request {
	var jsonData []byte
	var err error

	if data != nil {
		jsonData, err = json.Marshal(data)
		assert.NoError(t, err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	assert.NoError(t, err)

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

func makeDELETERequest(t *testing.T, url string) *http.Request {
	req, err := http.NewRequest("DELETE", url, http.NoBody)
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
	targetService := services.NewTarget(store.Target(), store.Mission())
	mockLogger := &MockLogger{}

	service := NewImplService(targetService)
	handler := NewHandler(service, mockLogger)

	router := gin.New()
	router.Use(middleware.GlobalErrorHandler())
	v1 := router.Group("/v1")
	targets := v1.Group("/missions/:id/targets")
	{
		targets.POST("", handler.Add)
		targets.PUT("/:tid/notes", handler.UpdateNotes)
		targets.PUT("/:tid/complete", handler.MarkComplete)
		targets.DELETE("/:tid", handler.Delete)
	}

	return router, service
}

func TestTargetController_Add(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should add target to existing mission", func(t *testing.T) {
		target := createTestTarget("New Target", "Spain", "New target notes")
		req := makePOSTRequest(t, TargetsBaseURL+"/1/targets", target)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response models.Target
		assertJSONResponse(t, w, http.StatusCreated, &response)
		assert.Equal(t, "New Target", response.Name)
		assert.NotZero(t, response.ID)
	})

	t.Run("should fail for completed mission", func(t *testing.T) {
		target := createTestTarget("Test Target", "Test", "Test notes")
		req := makePOSTRequest(t, TargetsBaseURL+"/2/targets", target)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail for non-existing mission", func(t *testing.T) {
		target := createTestTarget("Test Target", "Test", "Test notes")
		req := makePOSTRequest(t, TargetsBaseURL+"/999/targets", target)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusNotFound, nil)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		req, err := http.NewRequest("POST", TargetsBaseURL+"/1/targets", bytes.NewBuffer([]byte(InvalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail with invalid mission ID", func(t *testing.T) {
		target := createTestTarget("Test Target", "Test", "Test notes")
		req := makePOSTRequest(t, TargetsBaseURL+"/"+InvalidMissionIDStr+"/targets", target)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail with empty target name", func(t *testing.T) {
		target := createTestTarget("", "Test", "Test notes") // Empty name
		req := makePOSTRequest(t, TargetsBaseURL+"/1/targets", target)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail with very long notes", func(t *testing.T) {
		longNotes := strings.Repeat("a", 501) // Notes longer than 500 characters
		target := createTestTarget("Test Target", "Test", longNotes)
		req := makePOSTRequest(t, TargetsBaseURL+"/1/targets", target)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})
}

func TestTargetController_UpdateNotes(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should update target notes", func(t *testing.T) {
		updateReq := UpdateNotesRequest{Notes: "Updated notes for target"}
		req := makePUTRequest(t, TargetsBaseURL+"/1/targets/1/notes", updateReq)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusOK, nil)
	})

	t.Run("should fail for completed mission", func(t *testing.T) {
		updateReq := UpdateNotesRequest{Notes: "Test notes"}
		jsonData, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/missions/2/targets/3/notes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // ErrMissionComplete has status 400
	})

	t.Run("should fail for non-existing target", func(t *testing.T) {
		updateReq := UpdateNotesRequest{Notes: "Test notes"}
		jsonData, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/missions/1/targets/999/notes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/v1/missions/1/targets/1/notes", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid mission ID", func(t *testing.T) {
		updateReq := UpdateNotesRequest{Notes: "Test notes"}
		jsonData, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/missions/invalid/targets/1/notes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid target ID", func(t *testing.T) {
		updateReq := UpdateNotesRequest{Notes: "Test notes"}
		jsonData, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/missions/1/targets/invalid/notes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTargetController_MarkComplete(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should mark target as complete", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/v1/missions/1/targets/2/complete", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should fail for non-existing target", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/v1/missions/1/targets/999/complete", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should fail with invalid target ID", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/v1/missions/1/targets/invalid/complete", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTargetController_Delete(t *testing.T) {
	router, service := setupTestRouter()

	t.Run("should delete incomplete target", func(t *testing.T) {
		// First add a target to delete
		target := &models.Target{
			Name:     "To Delete",
			Country:  "Test",
			Notes:    "Test",
			Complete: false,
		}
		err := service._targetContext.Add(context.TODO(), 1, target)
		assert.NoError(t, err)

		req := makeDELETERequest(t, TargetsBaseURL+"/1/targets/"+strconv.FormatUint(uint64(target.ID), 10))

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("should fail for completed target", func(t *testing.T) {
		req := makeDELETERequest(t, TargetsBaseURL+"/2/targets/3") // Target 3 is complete

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // ErrTargetComplete has status 400
	})

	t.Run("should fail for non-existing target", func(t *testing.T) {
		req := makeDELETERequest(t, TargetsBaseURL+"/1/targets/999")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should fail with invalid mission ID", func(t *testing.T) {
		req := makeDELETERequest(t, TargetsBaseURL+"/"+InvalidMissionIDStr+"/targets/1")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid target ID", func(t *testing.T) {
		req := makeDELETERequest(t, TargetsBaseURL+"/1/targets/"+InvalidTargetIDStr)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
