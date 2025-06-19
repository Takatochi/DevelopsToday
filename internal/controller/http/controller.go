package http

import (
	"DevelopsToday/config"
	v1 "DevelopsToday/internal/controller/http/v1"
	"DevelopsToday/internal/controller/http/v1/cat"
	"DevelopsToday/internal/controller/http/v1/mission"
	"DevelopsToday/internal/controller/http/v1/target"
	"DevelopsToday/internal/repo"
	"DevelopsToday/internal/services"
	// Swagger documentation
	_ "DevelopsToday/docs"
	"DevelopsToday/internal/controller/http/middleware"
	"DevelopsToday/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Spy Cat Agency API
// @version		1.0
// @description	API for managing spy cats, missions, and targets
// @host			localhost:8080
// @BasePath		/v1
func NewV1Controller(engine *gin.Engine, store repo.Store, cfg *config.Config, l logger.Interface) {
	// Middleware
	engine.Use(middleware.LoggerMiddleware(l))
	engine.Use(middleware.RecoveryMiddleware(l))

	//// Swagger
	if cfg.Swagger.Enabled {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	//Repositories
	catRepo := store.Cat()
	missionRepo := store.Mission()
	targetRepo := store.Target()

	//Services
	catHandlerService := cat.NewImplService(
		services.NewBreed(),
		services.NewCat(catRepo),
	)

	missionHandlerService := mission.NewImplService(
		services.NewMission(missionRepo),
	)

	targetHandlerService := target.NewImplService(
		services.NewTarget(targetRepo, missionRepo),
	)

	// API v1 group
	v1Group := engine.Group("/v1")
	{
		v1.NewSpyCatsRoutes(v1Group, catHandlerService, l)
		v1.NewMissionsRoutes(v1Group, missionHandlerService, l)
		v1.NewTargetsRoutes(v1Group, targetHandlerService, l)
	}
}
