package http

import (
	"DevelopsToday/config"
	// Swagger documentation
	//_ "DevelopsToday/docs"
	"DevelopsToday/internal/controller/http/middleware"
	"DevelopsToday/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewV1Controller - inits controller for v1 API.
// Swagger spec:
// @title       API V1
// @description API V1 for DevelopsToday application
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
// @in header
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
		v1Group.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}
}
