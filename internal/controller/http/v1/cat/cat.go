package cat

import (
	"net/http"
	"strconv"

	"DevelopsToday/internal/models"
	"DevelopsToday/internal/services"

	"github.com/gin-gonic/gin"
)

type Service struct {
	_validator  services.Validator
	_catContext services.CatContext
}

func NewImplService(validator services.Validator, catCtx services.CatContext) *Service {
	return &Service{
		_validator:  validator,
		_catContext: catCtx,
	}
}

type Handler struct {
	Service *Service
}

// Create godoc
//
//	@Summary		Create a new cat
//	@Description	Create a new cat in the system
//	@Tags			cats
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			cat	body		models.Cat	true	"Cat info"
//	@Success		201	{object}	models.Cat
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Router			/cats [post]
func (h *Handler) Create(ctx *gin.Context) {
	var newCat models.Cat
	if err := ctx.ShouldBindJSON(&newCat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if !h.Service._validator.IsValid(newCat.Breed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid breed"})
		return
	}

	if err := h.Service._catContext.Create(ctx, &newCat); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cat"})
		return
	}

	ctx.JSON(http.StatusCreated, newCat)
}

// List godoc
//
//	@Summary		List all cats
//	@Description	Get list of all cats
//	@Tags			cats
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}	models.Cat
//	@Failure		401	{object}	map[string]interface{}
//	@Router			/cats [get]
func (h *Handler) List(ctx *gin.Context) {
	cats, err := h.Service._catContext.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cats"})
		return
	}
	ctx.JSON(http.StatusOK, cats)
}

//	@Des
//
// GetByID godoc
//
//	@Summary		Get cat by ID
//	@Description	Get cat details by ID
//	@Tags			cats
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Cat ID"
//	@Success		200	{object}	models.Cat
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/cats/{id} [get]
func (h *Handler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	cat, err := h.Service._catContext.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	ctx.JSON(http.StatusOK, cat)
}

// UpdateSalary godoc
//
//	@Summary		Update cat salary
//	@Description	Update salary for a specific cat
//	@Tags			cats
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int		true	"Cat ID"
//	@Param			salary	body		float64	true	"New salary"
//	@Success		200		{object}	models.Cat
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]interface{}
//	@Router			/cats/{id}/salary [patch]
func (h *Handler) UpdateSalary(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var body struct {
		Salary float64 `json:"salary"`
	}
	if bindErr := ctx.ShouldBindJSON(&body); bindErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if updateErr := h.Service._catContext.UpdateSalary(ctx, uint(id), body.Salary); updateErr != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	cat, err := h.Service._catContext.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Salary updated, but cat fetch failed"})
		return
	}

	ctx.JSON(http.StatusOK, cat)
}

// Delete godoc
//
//	@Summary		Delete cat
//	@Description	Delete cat by ID
//	@Tags			cats
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Cat ID"
//	@Success		204	"No Content"
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Router			/cats/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.Service._catContext.DeleteByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
