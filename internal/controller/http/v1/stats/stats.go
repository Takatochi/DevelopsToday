package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"

	mdware "DevelopsToday/internal/controller/http/middleware"
	"DevelopsToday/internal/services"
	"DevelopsToday/pkg/logger"
)

// Service wraps stats context
type Service struct {
	_statsContext services.StatsContext
}

// NewImplService creates a new stats service
func NewImplService(statsContext services.StatsContext) *Service {
	return &Service{
		_statsContext: statsContext,
	}
}

// Handler handles stats HTTP requests
type Handler struct {
	service *Service
	logger  logger.Interface
}

// NewHandler creates a new stats handler
func NewHandler(service *Service, logger logger.Interface) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// GetDashboard godoc
// @Summary Get dashboard statistics
// @Description Get comprehensive dashboard statistics including cats, missions, and targets data
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} services.DashboardStats
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /stats/dashboard [get]
func (h *Handler) GetDashboard(ctx *gin.Context) {
	stats, err := h.service._statsContext.GetDashboard(ctx)
	if err != nil {
		_ = ctx.Error(mdware.ErrInternalError)
		h.logger.Error("Failed to get dashboard stats: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, stats)
}
