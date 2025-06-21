package bulk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	mdware "DevelopsToday/internal/controller/http/middleware"
	"DevelopsToday/internal/models"
	"DevelopsToday/internal/services"
	"DevelopsToday/pkg/logger"
)

// Service wraps bulk context
type Service struct {
	_bulkContext services.BulkContext
}

// NewImplService creates a new bulk service
func NewImplService(bulkContext services.BulkContext) *Service {
	return &Service{
		_bulkContext: bulkContext,
	}
}

// Handler handles bulk HTTP requests
type Handler struct {
	service   *Service
	logger    logger.Interface
	validator *validator.Validate
}

// NewHandler creates a new bulk handler
func NewHandler(service *Service, logger logger.Interface) *Handler {
	return &Handler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// BulkSalaryUpdateRequest represents bulk salary update request
type BulkSalaryUpdateRequest struct {
	Updates []services.SalaryUpdate `json:"updates" validate:"required,min=1,max=100"`
}

// BulkCreateCatsRequest represents bulk cat creation request
type BulkCreateCatsRequest struct {
	Cats []*models.Cat `json:"cats" validate:"required,min=1,max=50"`
}

// BulkUpdateSalary godoc
// @Summary Bulk update cat salaries
// @Description Update multiple cat salaries in parallel using worker pool
// @Tags bulk
// @Accept json
// @Produce json
// @Param request body BulkSalaryUpdateRequest true "Bulk salary update request"
// @Success 200 {object} services.BulkResult
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /bulk/cats/salary [put]
func (h *Handler) BulkUpdateSalary(ctx *gin.Context) {
	var req BulkSalaryUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	// Validate struct fields
	if err := h.validator.Struct(&req); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	// Validate each update
	for _, update := range req.Updates {
		if err := h.validator.Struct(&update); err != nil {
			_ = ctx.Error(mdware.ErrBadRequest)
			return
		}
	}

	result, err := h.service._bulkContext.BulkUpdateSalary(ctx, req.Updates)
	if err != nil {
		_ = ctx.Error(mdware.ErrInternalError)
		h.logger.Error("Failed to bulk update salaries: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// BulkCreateCats godoc
// @Summary Bulk create cats
// @Description Create multiple cats in parallel using worker pool
// @Tags bulk
// @Accept json
// @Produce json
// @Param request body BulkCreateCatsRequest true "Bulk cat creation request"
// @Success 201 {object} services.BulkResult
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /bulk/cats [post]
func (h *Handler) BulkCreateCats(ctx *gin.Context) {
	var req BulkCreateCatsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	// Validate struct fields
	if err := h.validator.Struct(&req); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	// Validate each cat
	for _, cat := range req.Cats {
		if err := h.validator.Struct(cat); err != nil {
			_ = ctx.Error(mdware.ErrBadRequest)
			return
		}
	}

	result, err := h.service._bulkContext.BulkCreateCats(ctx, req.Cats)
	if err != nil {
		_ = ctx.Error(mdware.ErrInternalError)
		h.logger.Error("Failed to bulk create cats: %v", err)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
