package mission

import (
	"DevelopsToday/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	missions []models.Mission
}

var missionIDCounter uint = 1
var targetIDCounter uint = 1

func (h *Handler) Create(ctx *gin.Context) {
	var input struct {
		Targets []models.Target `json:"targets"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil || len(input.Targets) < 1 || len(input.Targets) > 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid targets"})
		return
	}

	mission := models.Mission{ID: missionIDCounter}
	for i := range input.Targets {
		input.Targets[i].ID = targetIDCounter
		targetIDCounter++
		mission.Targets = append(mission.Targets, input.Targets[i])
	}

	missionIDCounter++
	h.missions = append(h.missions, mission)
	ctx.JSON(http.StatusCreated, mission)
}

func (h *Handler) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.missions)
}

func (h *Handler) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for _, m := range h.missions {
		if int(m.ID) == id {
			ctx.JSON(http.StatusOK, m)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
}

func (h *Handler) AssignCat(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var body struct {
		CatID uint `json:"cat_id"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	for i := range h.missions {
		if int(h.missions[i].ID) == id {
			h.missions[i].CatID = &body.CatID
			ctx.JSON(http.StatusOK, h.missions[i])
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
}

func (h *Handler) MarkComplete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for i := range h.missions {
		m := &h.missions[i]
		if int(m.ID) == id {
			for _, t := range m.Targets {
				if !t.Complete {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Not all targets are completed"})
					return
				}
			}
			m.Complete = true
			ctx.JSON(http.StatusOK, m)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, m := range h.missions {
		if int(m.ID) == id {
			if m.CatID != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Mission has assigned cat"})
				return
			}
			h.missions = append(h.missions[:i], h.missions[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Mission not found"})
}
