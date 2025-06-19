package cat

import (
	"DevelopsToday/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cats []models.Cat
}

var idCounter uint = 1

func (h *Handler) Create(ctx *gin.Context) {
	var newCat models.Cat
	if err := ctx.ShouldBindJSON(&newCat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	//if !IsValidBreed(newCat.Breed) {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid breed"})
	//	return
	//}

	newCat.ID = idCounter
	idCounter++
	h.cats = append(h.cats, newCat)

	ctx.JSON(http.StatusCreated, newCat)
}

func (h *Handler) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.cats)
}

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
