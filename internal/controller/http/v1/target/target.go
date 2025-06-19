package target

import (
	"DevelopsToday/internal/models"
	"DevelopsToday/internal/services"
	"net/http"
	"strconv"

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
	missions []models.Mission
	Service  *Service
}

// UpdateNotesRequest represents request body for updating target notes
type UpdateNotesRequest struct {
	Notes string `json:"notes" example:"Target usually visits gym at 6 PM"`
}

var targetIDCounter uint = 1

// Add godoc
//
//	@Summary		Add target to mission
//	@Description	Add a new target to an existing mission
//	@Tags			targets
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"Mission ID"
//	@Param			input	body		models.Target	true	"Target info"
//	@Success		201		{object}	models.Target
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Router			/missions/{id}/targets [post]
func (h *Handler) Add(ctx *gin.Context) {
	mid, _ := strconv.Atoi(ctx.Param("id"))
	var input models.Target
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target"})
		return
	}

	for i := range h.missions {
		if int(h.missions[i].ID) == mid {
			if h.missions[i].Complete {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Mission is completed"})
				return
			}
			input.ID = targetIDCounter
			targetIDCounter++
			h.missions[i].Targets = append(h.missions[i].Targets, input)
			ctx.JSON(http.StatusCreated, input)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
}

// UpdateNotes godoc
//
//	@Summary		Update target notes
//	@Description	Update notes for a specific target
//	@Tags			targets
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Mission ID"
//	@Param			tid		path		int					true	"Target ID"
//	@Param			input	body		UpdateNotesRequest	true	"Target notes"
//	@Success		200		{object}	models.Target
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		403		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Router			/missions/{id}/targets/{tid}/notes [patch]
func (h *Handler) UpdateNotes(ctx *gin.Context) {
	mid, _ := strconv.Atoi(ctx.Param("id"))
	tid, _ := strconv.Atoi(ctx.Param("tid"))
	var body UpdateNotesRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notes"})
		return
	}

	for i := range h.missions {
		if int(h.missions[i].ID) == mid {
			if h.missions[i].Complete {
				ctx.JSON(http.StatusForbidden, gin.H{"error": "Mission is completed"})
				return
			}
			for j := range h.missions[i].Targets {
				if int(h.missions[i].Targets[j].ID) == tid {
					if h.missions[i].Targets[j].Complete {
						ctx.JSON(http.StatusForbidden, gin.H{"error": "Target is completed"})
						return
					}
					h.missions[i].Targets[j].Notes = body.Notes
					ctx.JSON(http.StatusOK, h.missions[i].Targets[j])
					return
				}
			}
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Target or Mission not found"})
}

// MarkComplete godoc
//
//	@Summary		Mark target as complete
//	@Description	Mark a specific target as completed
//	@Tags			targets
//	@Produce		json
//	@Param			id	path		int	true	"Mission ID"
//	@Param			tid	path		int	true	"Target ID"
//	@Success		200	{object}	models.Target
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id}/targets/{tid}/complete [post]
func (h *Handler) MarkComplete(ctx *gin.Context) {
	mid, _ := strconv.Atoi(ctx.Param("id"))
	tid, _ := strconv.Atoi(ctx.Param("tid"))
	for i := range h.missions {
		if int(h.missions[i].ID) == mid {
			for j := range h.missions[i].Targets {
				if int(h.missions[i].Targets[j].ID) == tid {
					h.missions[i].Targets[j].Complete = true
					ctx.JSON(http.StatusOK, h.missions[i].Targets[j])
					return
				}
			}
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
}

// Delete godoc
//
//	@Summary		Delete target
//	@Description	Delete a target from mission
//	@Tags			targets
//	@Produce		json
//	@Param			id	path	int	true	"Mission ID"
//	@Param			tid	path	int	true	"Target ID"
//	@Success		204	"No Content"
//	@Failure		403	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id}/targets/{tid} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	mid, _ := strconv.Atoi(ctx.Param("id"))
	tid, _ := strconv.Atoi(ctx.Param("tid"))
	for i := range h.missions {
		if int(h.missions[i].ID) == mid {
			for j, t := range h.missions[i].Targets {
				if int(t.ID) == tid {
					if t.Complete {
						ctx.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete completed target"})
						return
					}
					h.missions[i].Targets = append(h.missions[i].Targets[:j], h.missions[i].Targets[j+1:]...)
					ctx.Status(http.StatusNoContent)
					return
				}
			}
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Target not found"})
}
