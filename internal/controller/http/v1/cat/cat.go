package cat

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	mdware "DevelopsToday/internal/controller/http/middleware"
	"DevelopsToday/internal/models"
	"DevelopsToday/internal/services"
	"DevelopsToday/pkg/logger"
)

const (
	_recordNotFound = "record not found"
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
	service *Service
	logger  logger.Interface
}

func NewHandler(service *Service, logger logger.Interface) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
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
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	if !h.service._validator.IsValid(newCat.Breed) {
		_ = ctx.Error(mdware.ErrInvalidBreed)
		return
	}

	if err := h.service._catContext.Create(ctx, &newCat); err != nil {
		h.logger.Error("Failed to create cat: %v", err)
		_ = ctx.Error(mdware.ErrInternalError)
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
	cats, err := h.service._catContext.GetAll(ctx)
	if err != nil {
		_ = ctx.Error(mdware.ErrInternalError)
		h.logger.Error("Failed to list cats: %v", err)
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
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	cat, err := h.service._catContext.GetByID(ctx, uint(id))
	if err != nil {
		if err.Error() == _recordNotFound {
			_ = ctx.Error(mdware.ErrCatNotFound)
		} else {
			_ = ctx.Error(mdware.ErrInternalError)
		}
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
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	var body struct {
		Salary float64 `json:"salary"`
	}
	if bindErr := ctx.ShouldBindJSON(&body); bindErr != nil {
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	if updateErr := h.service._catContext.UpdateSalary(ctx, uint(id), body.Salary); updateErr != nil {
		if updateErr.Error() == _recordNotFound {
			_ = ctx.Error(mdware.ErrCatNotFound)
		} else {
			_ = ctx.Error(mdware.ErrInternalError)
		}
		h.logger.Error("Failed to update cat salary: %v", updateErr)
		return
	}

	cat, err := h.service._catContext.GetByID(ctx, uint(id))
	if err != nil {
		if err.Error() == _recordNotFound {
			_ = ctx.Error(mdware.ErrCatNotFound)
		} else {
			_ = ctx.Error(mdware.ErrInternalError)
		}
		h.logger.Error("Failed to fetch cat after salary update: %v", err)
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
		_ = ctx.Error(mdware.ErrBadRequest)
		return
	}

	err = h.service._catContext.DeleteByID(ctx, uint(id))
	if err != nil {
		if err.Error() == _recordNotFound {
			_ = ctx.Error(mdware.ErrCatNotFound)
		} else {
			_ = ctx.Error(mdware.ErrInternalError)
		}
		h.logger.Error("Failed to delete cat: %v", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
