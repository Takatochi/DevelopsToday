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

	missions := apiV1Group.Group("/missions")
	{
		missions.POST("/add", missionHandler.Create)
		missions.GET("", missionHandler.List)
		missions.GET("/:id", missionHandler.GetByID)
		missions.PUT("/:id/complete", missionHandler.MarkComplete)
		missions.PUT("/:id/assign", missionHandler.AssignCat)
		missions.DELETE("/:id", missionHandler.Delete)
	}

}
func NewTargetsRoutes(apiV1Group *gin.RouterGroup, l logger.Interface) {

	targets := apiV1Group.Group("/missions/:id/targets")
	{
		targets.POST("/add", targetHandler.Add)
		targets.DELETE("/:tid", targetHandler.Delete)
		targets.PUT("/:tid", targetHandler.UpdateNotes)
		targets.PUT("/:tid/complete", targetHandler.MarkComplete)
	}

}
