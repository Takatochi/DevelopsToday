package mission

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	mdware "DevelopsToday/internal/controller/http/middleware"
	"DevelopsToday/internal/models"
	"DevelopsToday/internal/services"
	"DevelopsToday/pkg/logger"
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
	service   *Service
	logger    logger.Interface
	validator *validator.Validate
}

func NewHandler(service *Service, logger logger.Interface) *Handler {
	return &Handler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// CreateRequest represents the request body for creating a mission
type CreateRequest struct {
	Targets []models.Target `json:"targets" validate:"required,min=1,max=3,dive"`
}

// AssignCatRequest represents the request body for assigning a cat
type AssignCatRequest struct {
	CatID uint `json:"cat_id" validate:"required,min=1"`
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
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	// Validate struct fields
	if err := h.validator.Struct(&input); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	mission := &models.Mission{
		Targets: input.Targets,
	}

	if err := h.service._missionContext.Create(ctx, mission); err != nil {
		_ = ctx.Error(mdware.ErrInternalError)
		h.logger.Error("Failed to create mission: %v", err)
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
		_ = ctx.Error(mdware.ErrInternalError)
		h.logger.Error("Failed to list missions: %v", err)
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	mission, err := h.service._missionContext.GetByID(ctx, uint(id))
	if err != nil {
		_ = ctx.Error(mdware.ErrMissionNotFound)
		h.logger.Error("Failed to get mission by ID: %v", err)
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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	var body AssignCatRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	// Validate struct fields
	if err := h.validator.Struct(&body); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	if err := h.service._missionContext.AssignCat(ctx, uint(id), body.CatID); err != nil {
		_ = ctx.Error(mdware.ErrInternalError)
		h.logger.Error("Failed to assign cat: %v", err)
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
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	// Викликаємо сервісний метод
	err = h.service._missionContext.MarkComplete(ctx, uint(id))
	if err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		h.logger.Error("Failed to mark mission complete: %v", err)
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
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	err = h.service._missionContext.DeleteByID(ctx, uint(id))
	if err != nil {
		if err.Error() == "cannot delete assigned mission" {
			_ = ctx.Error(mdware.ErrBadRequest)
		} else {
			_ = ctx.Error(mdware.ErrMissionNotFound)
		}
		h.logger.Error("Failed to delete mission: %v", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
