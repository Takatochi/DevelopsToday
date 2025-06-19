package target

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
	targetService := services.NewTarget(store.Target(), store.Mission())

	service := NewImplService(targetService)
	handler := &Handler{Service: service}

	router := gin.New()
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
		target := models.Target{
			Name:     "New Target",
			Country:  "Spain",
			Notes:    "New target notes",
			Complete: false,
		}

		jsonData, _ := json.Marshal(target)
		req, _ := http.NewRequest("POST", "/v1/missions/1/targets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Target
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "New Target", response.Name)
		assert.NotZero(t, response.ID)
	})

	t.Run("should fail for completed mission", func(t *testing.T) {
		target := models.Target{
			Name:     "Test Target",
			Country:  "Test",
			Notes:    "Test",
			Complete: false,
		}

		jsonData, _ := json.Marshal(target)
		req, _ := http.NewRequest("POST", "/v1/missions/2/targets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail for non-existing mission", func(t *testing.T) {
		target := models.Target{
			Name:     "Test Target",
			Country:  "Test",
			Notes:    "Test",
			Complete: false,
		}

		jsonData, _ := json.Marshal(target)
		req, _ := http.NewRequest("POST", "/v1/missions/999/targets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should fail with invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/v1/missions/1/targets", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid mission ID", func(t *testing.T) {
		target := models.Target{
			Name:     "Test Target",
			Country:  "Test",
			Notes:    "Test",
			Complete: false,
		}

		jsonData, _ := json.Marshal(target)
		req, _ := http.NewRequest("POST", "/v1/missions/invalid/targets", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTargetController_UpdateNotes(t *testing.T) {
	router, _ := setupTestRouter()

	t.Run("should update target notes", func(t *testing.T) {
		updateReq := UpdateNotesRequest{Notes: "Updated notes for target"}
		jsonData, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/missions/1/targets/1/notes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should fail for completed mission", func(t *testing.T) {
		updateReq := UpdateNotesRequest{Notes: "Test notes"}
		jsonData, _ := json.Marshal(updateReq)

		req, _ := http.NewRequest("PUT", "/v1/missions/2/targets/3/notes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
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

		req, _ := http.NewRequest("DELETE", "/v1/missions/1/targets/"+strconv.FormatUint(uint64(target.ID), 10), http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("should fail for completed target", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/missions/2/targets/3", http.NoBody) // Target 3 is complete

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("should fail for non-existing target", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/missions/1/targets/999", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should fail with invalid mission ID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/missions/invalid/targets/1", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should fail with invalid target ID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/v1/missions/1/targets/invalid", http.NoBody)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
