package target

import (
	"DevelopsToday/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	missions []models.Mission
}

var targetIDCounter uint = 1

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

func (h *Handler) UpdateNotes(ctx *gin.Context) {
	mid, _ := strconv.Atoi(ctx.Param("id"))
	tid, _ := strconv.Atoi(ctx.Param("tid"))
	var body struct {
		Notes string `json:"notes"`
	}
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
