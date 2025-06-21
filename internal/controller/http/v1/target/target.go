package target

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
	_targetContext services.TargetContext
}

func NewImplService(targetContext services.TargetContext) *Service {
	return &Service{
		_targetContext: targetContext,
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

// UpdateNotesRequest represents request body for updating target notes
type UpdateNotesRequest struct {
	Notes string `json:"notes" validate:"required,min=1,max=500" example:"Target usually visits gym at 6 PM"`
}

// Add godoc
//
//	@Summary		Add target to mission
//	@Description	Add a new target to an existing mission
//	@Tags			targets
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int				true	"Mission ID"
//	@Param			input	body		models.Target	true	"Target info"
//	@Success		201		{object}	models.Target
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Router			/missions/{id}/targets [post]
func (h *Handler) Add(ctx *gin.Context) {
	mid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	var input models.Target
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	// Validate struct fields
	if err := h.validator.Struct(&input); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	if err := h.service._targetContext.Add(ctx, uint(mid), &input); err != nil {
		if err.Error() == "cannot add target to completed mission" {
			_ = ctx.Error(mdware.ErrMissionComplete)
		} else {
			_ = ctx.Error(mdware.ErrMissionNotFound)
		}
		h.logger.Error("Failed to add target: %v", err)
		return
	}

	ctx.JSON(http.StatusCreated, input)
}

// UpdateNotes godoc
//
//	@Summary		Update target notes
//	@Description	Update notes for a specific target
//	@Tags			targets
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int					true	"Mission ID"
//	@Param			tid		path		int					true	"Target ID"
//	@Param			input	body		UpdateNotesRequest	true	"Target notes"
//	@Success		200		{object}	models.Target
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		403		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Router			/missions/{id}/targets/{tid}/notes [patch]
func (h *Handler) UpdateNotes(ctx *gin.Context) {
	mid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	tid, err := strconv.Atoi(ctx.Param("tid"))
	if err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	var body UpdateNotesRequest
	if err = ctx.ShouldBindJSON(&body); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	// Validate struct fields
	if err := h.validator.Struct(&body); err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	if err = h.service._targetContext.UpdateNotes(ctx, uint(mid), uint(tid), body.Notes); err != nil {
		if err.Error() == "mission is completed" {
			_ = ctx.Error(mdware.ErrMissionComplete)
		} else if err.Error() == "target is completed" {
			_ = ctx.Error(mdware.ErrTargetComplete)
		} else {
			_ = ctx.Error(mdware.ErrTargetNotFound)
		}
		h.logger.Error("Failed to update target notes: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Notes updated successfully"})
}

// MarkComplete godoc
//
//	@Summary		Mark target as complete
//	@Description	Mark a specific target as completed
//	@Tags			targets
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Mission ID"
//	@Param			tid	path		int	true	"Target ID"
//	@Success		200	{object}	models.Target
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id}/targets/{tid}/complete [post]
func (h *Handler) MarkComplete(ctx *gin.Context) {
	tid, err := strconv.Atoi(ctx.Param("tid"))
	if err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	if err := h.service._targetContext.MarkComplete(ctx, uint(tid)); err != nil {
		_ = ctx.Error(mdware.ErrTargetNotFound)
		h.logger.Error("Failed to mark target complete: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Target marked as complete"})
}

// Delete godoc
//
//	@Summary		Delete target
//	@Description	Delete a target from mission
//	@Tags			targets
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Mission ID"
//	@Param			tid	path	int	true	"Target ID"
//	@Success		204	"No Content"
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		403	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id}/targets/{tid} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	mid, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	tid, err := strconv.Atoi(ctx.Param("tid"))
	if err != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	if err = h.service._targetContext.DeleteByID(ctx, uint(mid), uint(tid)); err != nil {
		if err.Error() == "cannot delete completed target" {
			_ = ctx.Error(mdware.ErrTargetComplete)
		} else {
			_ = ctx.Error(mdware.ErrTargetNotFound)
		}
		h.logger.Error("Failed to delete target: %v", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
