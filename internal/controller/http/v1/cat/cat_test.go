package cat

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
	// Test cat data
	TestCatID     = uint(1)
	TestCatName   = "Whiskers"
	InvalidCatID  = uint(999)
	TestSalary    = float64(1000)
	UpdatedSalary = float64(2000)

	// Test URLs
	CatsBaseURL     = "/v1/cats"
	InvalidCatIDStr = "invalid"

	// Test JSON
	InvalidJSON = "invalid json"
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

// Helper functions for tests
func createTestCat(name, breed string, experience int, salary float64) models.Cat {
	return models.Cat{
		Name:       name,
		Experience: experience,
		Breed:      breed,
		Salary:     salary,
	}
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
		cat := createTestCat("TestCat", "Bengal", 3, TestSalary)
		req := makeJSONRequest(t, "POST", CatsBaseURL, cat)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response models.Cat
		assertJSONResponse(t, w, http.StatusCreated, &response)
		assert.Equal(t, "TestCat", response.Name)
		assert.NotZero(t, response.ID)
	})

	t.Run("should fail with invalid breed", func(t *testing.T) {
		cat := createTestCat("TestCat", "InvalidBreed", 2, TestSalary)
		req := makeJSONRequest(t, "POST", CatsBaseURL, cat)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		req, err := http.NewRequest("POST", CatsBaseURL, bytes.NewBuffer([]byte(InvalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail with empty name", func(t *testing.T) {
		cat := createTestCat("", "Bengal", 1, TestSalary) // Empty name
		req := makeJSONRequest(t, "POST", CatsBaseURL, cat)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail with negative salary", func(t *testing.T) {
		cat := createTestCat("TestCat", "Bengal", 4, -100) // Negative salary
		req := makeJSONRequest(t, "POST", CatsBaseURL, cat)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
	})

	t.Run("should fail with negative experience", func(t *testing.T) {
		cat := models.Cat{
			Name:       "TestCat",
			Experience: -1, // Negative experience
			Breed:      "Bengal",
			Salary:     TestSalary,
		}

		jsonData, err := json.Marshal(cat)
		assert.NoError(t, err)
		req, err := http.NewRequest("POST", CatsBaseURL, bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with very long name", func(t *testing.T) {
		longName := strings.Repeat("a", 101) // Name longer than 100 characters
		cat := models.Cat{
			Name:       longName,
			Experience: 5,
			Breed:      "Bengal",
			Salary:     TestSalary,
		}

		jsonData, err := json.Marshal(cat)
		assert.NoError(t, err)
		req, err := http.NewRequest("POST", CatsBaseURL, bytes.NewBuffer(jsonData))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCatController_List(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should return list of cats", func(t *testing.T) {
		req, err := http.NewRequest("GET", CatsBaseURL, http.NoBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Cat
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 5) // Should have at least 5 initial cats
	})
}

func TestCatController_GetByID(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should return cat by valid ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", CatsBaseURL+"/1", http.NoBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Cat
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, TestCatID, response.ID)
		assert.Equal(t, TestCatName, response.Name)
	})

	t.Run("should return 404 for non-existing cat", func(t *testing.T) {
		req, err := http.NewRequest("GET", CatsBaseURL+"/999", http.NoBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 400 for invalid ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", CatsBaseURL+"/"+InvalidCatIDStr, http.NoBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCatController_UpdateSalary(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should update cat salary", func(t *testing.T) {
		updateData := map[string]float64{"salary": UpdatedSalary}
		req := makeJSONRequest(t, "PUT", CatsBaseURL+"/1/salary", updateData)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusOK, nil)

		// Verify salary was updated
		req2, err := http.NewRequest("GET", CatsBaseURL+"/1", http.NoBody)
		assert.NoError(t, err)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		var response models.Cat
		err = json.Unmarshal(w2.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, UpdatedSalary, response.Salary)
	})

	t.Run("should return 404 for non-existing cat", func(t *testing.T) {
		updateData := map[string]float64{"salary": UpdatedSalary}
		req := makeJSONRequest(t, "PUT", CatsBaseURL+"/999/salary", updateData)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusNotFound, nil)
	})

	t.Run("should return 400 for invalid JSON", func(t *testing.T) {
		req, err := http.NewRequest("PUT", CatsBaseURL+"/1/salary", bytes.NewBuffer([]byte(InvalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with negative salary", func(t *testing.T) {
		updateData := map[string]float64{"salary": -100}
		req := makeJSONRequest(t, "PUT", CatsBaseURL+"/1/salary", updateData)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assertJSONResponse(t, w, http.StatusBadRequest, nil)
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

		req, err := http.NewRequest("DELETE", CatsBaseURL+"/"+strconv.Itoa(int(cat.ID)), http.NoBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify cat was deleted
		req2, err := http.NewRequest("GET", CatsBaseURL+"/"+strconv.Itoa(int(cat.ID)), http.NoBody)
		assert.NoError(t, err)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusNotFound, w2.Code)
	})

	t.Run("should return 404 for non-existing cat", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", CatsBaseURL+"/999", http.NoBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return 400 for invalid ID", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", CatsBaseURL+"/"+InvalidCatIDStr, http.NoBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
