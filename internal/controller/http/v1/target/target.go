package target

import (
	"net/http"
	"strconv"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/services"

	"github.com/gin-gonic/gin"
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
	Service *Service
}

// UpdateNotesRequest represents request body for updating target notes
type UpdateNotesRequest struct {
	Notes string `json:"notes" example:"Target usually visits gym at 6 PM"`
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	var input models.Target
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target"})
		return
	}

	if err := h.Service._targetContext.Add(ctx, uint(mid), &input); err != nil {
		if err.Error() == "cannot add target to completed mission" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
		}
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	tid, err := strconv.Atoi(ctx.Param("tid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	var body UpdateNotesRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notes"})
		return
	}

	if err := h.Service._targetContext.UpdateNotes(ctx, uint(mid), uint(tid), body.Notes); err != nil {
		if err.Error() == "mission is completed" || err.Error() == "target is completed" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Target or Mission not found"})
		}
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	if err := h.Service._targetContext.MarkComplete(ctx, uint(tid)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	tid, err := strconv.Atoi(ctx.Param("tid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	if err := h.Service._targetContext.DeleteByID(ctx, uint(mid), uint(tid)); err != nil {
		if err.Error() == "cannot delete completed target" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
