package http

import (
	"DevelopsToday/config"
	v1 "DevelopsToday/internal/controller/http/v1"

	// Swagger documentation
	_ "DevelopsToday/docs"
	"DevelopsToday/internal/controller/http/middleware"
	"DevelopsToday/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Spy Cat Agency API
//	@version		1.0
//	@description	API for managing spy cats, missions, and targets
//	@host			localhost:8080
//	@BasePath		/v1
func NewV1Controller(engine *gin.Engine, cfg *config.Config, l logger.Interface) {
	// Middleware
	engine.Use(middleware.LoggerMiddleware(l))
	engine.Use(middleware.RecoveryMiddleware(l))

	//// Swagger
	if cfg.Swagger.Enabled {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API v1 group
	v1Group := engine.Group("/v1")
	{
		v1.NewSpyCatsRoutes(v1Group, l)
		v1.NewMissionsRoutes(v1Group, l)
		v1.NewTargetsRoutes(v1Group, l)
	}
}
