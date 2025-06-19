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

// CreateRequest represents the request body for creating a mission
type CreateRequest struct {
	Targets []models.Target `json:"targets"`
}

// AssignCatRequest represents the request body for assigning a cat
type AssignCatRequest struct {
	CatID uint `json:"cat_id"`
}

// Create godoc
//	@Summary		Create a new mission
//	@Description	Create a new mission with 1-3 targets
//	@Tags			missions
//	@Accept			json
//	@Produce		json
//	@Param			input	body		CreateRequest	true	"Mission targets"
//	@Success		201		{object}	models.Mission
//	@Failure		400		{object}	map[string]interface{}
//	@Router			/missions [post]
func (h *Handler) Create(ctx *gin.Context) {
	var input CreateRequest
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

// List godoc
//	@Summary		List all missions
//	@Description	Get all missions
//	@Tags			missions
//	@Produce		json
//	@Success		200	{array}	models.Mission
//	@Router			/missions [get]
func (h *Handler) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.missions)
}

// GetByID godoc
//	@Summary		Get mission by ID
//	@Description	Get mission details by ID
//	@Tags			missions
//	@Produce		json
//	@Param			id	path		int	true	"Mission ID"
//	@Success		200	{object}	models.Mission
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id} [get]
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

// AssignCat godoc
//	@Summary		Assign cat to mission
//	@Description	Assign a cat to complete the mission
//	@Tags			missions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Mission ID"
//	@Param			input	body		AssignCatRequest	true	"Cat info"
//	@Success		200		{object}	models.Mission
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Router			/missions/{id}/assign [post]
func (h *Handler) AssignCat(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var body AssignCatRequest
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

// MarkComplete godoc
//	@Summary		Mark mission as complete
//	@Description	Mark mission as complete if all targets are completed
//	@Tags			missions
//	@Produce		json
//	@Param			id	path		int	true	"Mission ID"
//	@Success		200	{object}	models.Mission
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id}/complete [post]
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

// Delete godoc
//	@Summary		Delete mission
//	@Description	Delete mission if it has no assigned cat
//	@Tags			missions
//	@Produce		json
//	@Param			id	path	int	true	"Mission ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/missions/{id} [delete]
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