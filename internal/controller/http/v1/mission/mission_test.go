package mission

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo/mocks"
	"DevelopsToday/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *Service) {
	gin.SetMode(gin.TestMode)

	store := mocks.NewRepository()
	missionService := services.NewMission(store.Mission())

	service := NewImplService(missionService)
	handler := NewHandler(service, nil)

	router := gin.New()
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
			Targets: []models.Target{
				{Name: "Test Target 1", Country: "Ukraine", Notes: "Test notes 1", Complete: false},
				{Name: "Test Target 2", Country: "Poland", Notes: "Test notes 2", Complete: false},
			},
		}

		jsonData, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/v1/missions", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Mission
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotZero(t, response.ID)
		assert.Equal(t, 2, len(response.Targets))
	})

	t.Run("should fail with no targets", func(t *testing.T) {
		createReq := CreateRequest{
			Targets: []models.Target{},
		}

		jsonData, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/v1/missions", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with more than 3 targets", func(t *testing.T) {
		createReq := CreateRequest{
			Targets: []models.Target{
				{Name: "Target 1", Country: "Country 1", Notes: "Notes 1", Complete: false},
				{Name: "Target 2", Country: "Country 2", Notes: "Notes 2", Complete: false},
				{Name: "Target 3", Country: "Country 3", Notes: "Notes 3", Complete: false},
				{Name: "Target 4", Country: "Country 4", Notes: "Notes 4", Complete: false},
			},
		}

		jsonData, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", "/v1/missions", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/missions", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestMissionController_List(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should return list of missions", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/v1/missions", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Mission
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
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
