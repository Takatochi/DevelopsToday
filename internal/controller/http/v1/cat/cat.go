package cat

import (
	"DevelopsToday/internal/models"
	"DevelopsToday/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	services.Validator
}
type Handler struct {
	cats    []models.Cat
	Service Service
}

var idCounter uint = 1

// Create godoc
//
//	@Summary		Create a new
//	@Summary		Create a new cat
//	@Description	Create a new cat in the system
//	@Tags			cats
//	@Accept			json
//	@Produce		json
//	@Param			cat	body		models.Cat	true	"Cat info"
//	@Success		201	{object}	models.Cat
//	@Failure		400	{object}	map[string]interface{}
//	@Router			/cats [post]
func (h *Handler) Create(ctx *gin.Context) {
	var newCat models.Cat
	if err := ctx.ShouldBindJSON(&newCat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if !h.Service.IsValid(newCat.Breed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid breed"})
		return
	}

	newCat.ID = idCounter
	idCounter++
	h.cats = append(h.cats, newCat)

	ctx.JSON(http.StatusCreated, newCat)
}

// List godoc
//
//	@Summary		List all cats
//	@Description	Get list of all cats
//	@Tags			cats
//	@Produce		json
//	@Success		200	{array}	models.Cat
//	@Router			/cats [get]
func (h *Handler) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.cats)
}

//	@Des
//
// GetByID godoc
//
//	@Summary		Get cat by ID
//	@Description	Get cat details by ID
//	@Tags			cats
//	@Produce		json
//	@Param			id	path		int	true	"Cat ID"
//	@Success		200	{object}	models.Cat
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/cats/{id} [get]
func (h *Handler) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for _, cat := range h.cats {
		if int(cat.ID) == id {
			ctx.JSON(http.StatusOK, cat)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
}

// UpdateSalary godoc
//
//	@Summary		Update cat salary
//	@Description	Update salary for a specific cat
//	@Tags			cats
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Cat ID"
//	@Param			salary	body		float64	true	"New salary"
//	@Success		200		{object}	models.Cat
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Router			/cats/{id}/salary [patch]
func (h *Handler) UpdateSalary(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var body struct {
		Salary float64 `json:"salary"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	for i := range h.cats {
		if int(h.cats[i].ID) == id {
			h.cats[i].Salary = body.Salary
			ctx.JSON(http.StatusOK, h.cats[i])
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
}

// Delete godoc
//
//	@Summary		Delete cat
//	@Description	Delete cat by ID
//	@Tags			cats
//	@Produce		json
//	@Param			id	path	int	true	"Cat ID"
//	@Success		204	"No Content"
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/cats/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	for i, cat := range h.cats {
		if int(cat.ID) == id {
			h.cats = append(h.cats[:i], h.cats[i+1:]...)
			ctx.Status(http.StatusNoContent)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
}
