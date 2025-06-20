package v1

import (
	"DevelopsToday/internal/controller/http/v1/cat"
	"DevelopsToday/internal/controller/http/v1/mission"
	"DevelopsToday/internal/controller/http/v1/target"
	"DevelopsToday/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewSpyCatsRoutes(apiV1Group *gin.RouterGroup, service *cat.Service, l logger.Interface) {
	handler := &V1{
		cat: cat.NewHandler(service, l),
	}

	cats := apiV1Group.Group("/cats")
	cats.POST("", handler.cat.Create)
	cats.GET("", handler.cat.List)
	cats.GET("/:id", handler.cat.GetByID)
	cats.PUT("/:id/salary", handler.cat.UpdateSalary)
	cats.DELETE("/:id", handler.cat.Delete)
}
func NewMissionsRoutes(apiV1Group *gin.RouterGroup, service *mission.Service, l logger.Interface) {
	handler := &V1{
		mission: mission.NewHandler(service, l),
	}
	missions := apiV1Group.Group("/missions")
	missions.POST("", handler.mission.Create)
	missions.GET("", handler.mission.List)
	missions.GET("/:id", handler.mission.GetByID)
	missions.PUT("/:id/complete", handler.mission.MarkComplete)
	missions.PUT("/:id/assign", handler.mission.AssignCat)
	missions.DELETE("/:id", handler.mission.Delete)
}
func NewTargetsRoutes(apiV1Group *gin.RouterGroup, service *target.Service, l logger.Interface) {
	handler := &V1{
		target: target.NewHandler(service, l),
	}
	targets := apiV1Group.Group("/missions/:id/targets")
	targets.POST("", handler.target.Add)
	targets.DELETE("/:tid", handler.target.Delete)
	targets.PUT("/:tid/notes", handler.target.UpdateNotes)
	targets.PUT("/:tid/complete", handler.target.MarkComplete)
}
