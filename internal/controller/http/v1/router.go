package v1

import (
	"DevelopsToday/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewSpyCatsRoutes(apiV1Group *gin.RouterGroup, l logger.Interface) {
	handler := &V1{l: l}
	cats := apiV1Group.Group("/cats")
	{
		cats.POST("/add")
		cats.GET("", handler.cat.List)
		cats.GET("/:id", handler.cat.GetByID)
		cats.PUT("/:id/salary", handler.cat.UpdateSalary)
		cats.DELETE("/:id", handler.cat.Delete)
	}
}
func NewMissionsRoutes(apiV1Group *gin.RouterGroup, l logger.Interface) {
	handler := &V1{l: l}
	missions := apiV1Group.Group("/missions")
	{
		missions.POST("/add", handler.mission.Create)
		missions.GET("", handler.mission.List)
		missions.GET("/:id", handler.mission.GetByID)
		missions.PUT("/:id/complete", handler.mission.MarkComplete)
		missions.PUT("/:id/assign", handler.mission.AssignCat)
		missions.DELETE("/:id", handler.mission.Delete)
	}

}
func NewTargetsRoutes(apiV1Group *gin.RouterGroup, l logger.Interface) {
	handler := &V1{l: l}
	targets := apiV1Group.Group("/missions/:id/targets")
	{
		targets.POST("/add", handler.target.Add)
		targets.DELETE("/:tid", handler.target.Delete)
		targets.PUT("/:tid", handler.target.UpdateNotes)
		targets.PUT("/:tid/complete", handler.target.MarkComplete)
	}

}
