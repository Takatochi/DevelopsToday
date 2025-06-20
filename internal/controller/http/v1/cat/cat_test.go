package cat

import (
	"bytes"
	"context"
	"encoding/json"
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

// MockValidator для тестування
type MockValidator struct {
	validBreeds map[string]bool
}

func NewMockValidator() *MockValidator {
	return &MockValidator{
		validBreeds: map[string]bool{
			"Bengal":       true,
			"Siamese":      true,
			"Persian":      true,
			"Maine Coon":   true,
			"Russian Blue": true,
		},
	}
}

func (m *MockValidator) IsValid(breedName string) bool {
	return m.validBreeds[breedName]
}

// MockLogger для тестування
type MockLogger struct{}

func (m *MockLogger) Debug(message interface{}, args ...interface{}) {}
func (m *MockLogger) Info(message string, args ...interface{})       {}
func (m *MockLogger) Warn(message string, args ...interface{})       {}
func (m *MockLogger) Error(message interface{}, args ...interface{}) {}
func (m *MockLogger) Fatal(message interface{}, args ...interface{}) {}

func setupTestRouter() (*gin.Engine, *Service) {
	gin.SetMode(gin.TestMode)

	store := mocks.NewRepository()
	catService := services.NewCat(store.Cat())
	validator := NewMockValidator()
	mockLogger := &MockLogger{}

	service := NewImplService(validator, catService)
	handler := NewHandler(service, mockLogger)

	router := gin.New()
	router.Use(middleware.GlobalErrorHandler())
	v1 := router.Group("/v1")
	cats := v1.Group("/cats")
	{
		cats.POST("", handler.Create)
		cats.GET("", handler.List)
		cats.GET("/:id", handler.GetByID)
		cats.PUT("/:id/salary", handler.UpdateSalary)
		cats.DELETE("/:id", handler.Delete)
	}

	return router, service
}

func TestCatController_Create(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should create cat with valid data", func(t *testing.T) {
		cat := models.Cat{
			Name:       "TestCat",
			Experience: 5,
			Breed:      "Bengal",
			Salary:     1000,
		}

		jsonData, _ := json.Marshal(cat)
		req, _ := http.NewRequest("POST", "/v1/cats", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Cat
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "TestCat", response.Name)
		assert.NotZero(t, response.ID)
	})

	t.Run("should fail with invalid breed", func(t *testing.T) {
		cat := models.Cat{
			Name:       "TestCat",
			Experience: 5,
			Breed:      "InvalidBreed",
			Salary:     1000,
		}

		jsonData, _ := json.Marshal(cat)
		req, _ := http.NewRequest("POST", "/v1/cats", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/cats", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCatController_List(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should return list of cats", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/cats", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Cat
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 5) // Should have at least 5 initial cats
	})
}

func TestCatController_GetByID(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should return cat by valid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/cats/1", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Cat
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Whiskers", response.Name)
	})

	t.Run("should return 404 for non-existing cat", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/cats/999", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 400 for invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/cats/invalid", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCatController_UpdateSalary(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should update cat salary", func(t *testing.T) {
		updateData := map[string]float64{"salary": 2000}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", "/v1/cats/1/salary", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify salary was updated
		req2, _ := http.NewRequest("GET", "/v1/cats/1", http.NoBody)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		var response models.Cat
		err := json.Unmarshal(w2.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(2000), response.Salary)
	})

	t.Run("should return 404 for non-existing cat", func(t *testing.T) {
		updateData := map[string]float64{"salary": 2000}
		jsonData, _ := json.Marshal(updateData)

		req, _ := http.NewRequest("PUT", "/v1/cats/999/salary", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/v1/cats/1/salary", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCatController_Delete(t *testing.T) {
	router, service := setupTestRouter()

	t.Run("should delete existing cat", func(t *testing.T) {
		// First create a cat to delete
		cat := &models.Cat{
			Name:       "ToDelete",
			Experience: 1,
			Breed:      "Persian",
			Salary:     500,
		}
		err := service._catContext.Create(context.TODO(), cat)
		assert.NoError(t, err)

		req, _ := http.NewRequest("DELETE", "/v1/cats/"+strconv.Itoa(int(cat.ID)), http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify cat was deleted
		req2, _ := http.NewRequest("GET", "/v1/cats/"+strconv.Itoa(int(cat.ID)), http.NoBody)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusNotFound, w2.Code)
	})

	t.Run("should return 404 for non-existing cat", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/cats/999", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 400 for invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/cats/invalid", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
