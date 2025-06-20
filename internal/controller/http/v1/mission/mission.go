package mission

import (
	"DevelopsToday/pkg/logger"
	"net/http"
	"strconv"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/services"

	"github.com/gin-gonic/gin"
)

type Service struct {
	_missionContext services.MissionContext
}

func NewImplService(missionContext services.MissionContext) *Service {
	return &Service{
		_missionContext: missionContext,
	}
}

type Handler struct {
	service *Service
	logger  logger.Interface
}

func NewHandler(service *Service, logger logger.Interface) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// CreateRequest represents the request body for creating a mission
type CreateRequest struct {
	Targets []models.Target `json:"targets"`
}

// AssignCatRequest represents the request body for assigning a cat
type AssignCatRequest struct {
	CatID uint `json:"cat_id"`
}

// Create godoc
//
//	@Summary		Create a new mission
//	@Description	Create a new mission with 1-3 targets
//	@Tags			missions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		CreateRequest	true	"Mission targets"
//	@Success		201		{object}	models.Mission
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Router			/missions [post]
func (h *Handler) Create(ctx *gin.Context) {
	var input CreateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil || len(input.Targets) < 1 || len(input.Targets) > 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid targets"})
		return
	}

	mission := &models.Mission{
		Targets: input.Targets,
	}

	if err := h.service._missionContext.Create(ctx, mission); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mission"})
		return
	}

	ctx.JSON(http.StatusCreated, mission)
}

// List godoc
//
//	@Summary		List all missions
//	@Description	Get all missions
//	@Tags			missions
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	models.Mission
//	@Failure		401	{object}	map[string]interface{}
//	@Router			/missions [get]
func (h *Handler) List(ctx *gin.Context) {
	missions, err := h.service._missionContext.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list missions"})
		return
	}
	ctx.JSON(http.StatusOK, missions)
}

// GetByID godoc
//
//	@Summary		Get mission by ID
//	@Description	Get mission details by ID
//	@Tags			missions
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Mission ID"
//	@Success		200	{object}	models.Mission
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id} [get]
func (h *Handler) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	mission, err := h.service._missionContext.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		return
	}
	ctx.JSON(http.StatusOK, mission)
}

// AssignCat godoc
//
//	@Summary		Assign cat to mission
//	@Description	Assign a cat to complete the mission
//	@Tags			missions
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"Mission ID"
//	@Param			input	body		AssignCatRequest	true	"Cat info"
//	@Success		200		{object}	models.Mission
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Router			/missions/{id}/assign [post]
func (h *Handler) AssignCat(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var body AssignCatRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if err := h.service._missionContext.AssignCat(ctx, uint(id), body.CatID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign cat"})
		return
	}
	ctx.Status(http.StatusOK)
}

// MarkComplete godoc
//
//	@Summary		Mark mission as complete
//	@Description	Mark mission as complete if all targets are completed
//	@Tags			missions
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Mission ID"
//	@Success		200	{object}	models.Mission
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id}/complete [post]
func (h *Handler) MarkComplete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Викликаємо сервісний метод
	err = h.service._missionContext.MarkComplete(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Mission marked as complete"})
}

// Delete godoc
//
//	@Summary		Delete mission
//	@Description	Delete mission if it has no assigned cat
//	@Tags			missions
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Mission ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service._missionContext.DeleteByID(ctx, uint(id))
	if err != nil {
		if err.Error() == "cannot delete assigned mission" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
