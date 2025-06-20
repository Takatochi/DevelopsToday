package http

import (
	"DevelopsToday/config"
	v1 "DevelopsToday/internal/controller/http/v1"
	"DevelopsToday/internal/controller/http/v1/auth"
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

// NewV1Controller Swagger documentation
// @title Spy Cat Agency API
// @version 1.0
// @description REST API for managing spy cats, missions, and targets with JWT authentication
// @host			localhost:8080
// @BasePath		/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func NewV1Controller(
	engine *gin.Engine,
	store repo.Store,
	cfg *config.Config,
	l logger.Interface,
	jwtService *services.JWTService,
) {
	// Middleware
	engine.Use(middleware.LoggerMiddleware(l))
	engine.Use(middleware.RecoveryMiddleware(l))
	engine.Use(middleware.GlobalErrorHandler())

	// Health check endpoint
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	//// Swagger
	if cfg.Swagger.Enabled {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Repositories
	catRepo := store.Cat()
	missionRepo := store.Mission()
	targetRepo := store.Target()

	// Services
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

	// Auth handler
	authHandler := auth.NewHandler(store.User(), jwtService, l)

	// API v1 group
	v1Group := engine.Group("/v1")
	{
		// Auth routes (no authentication required)
		authGroup := v1Group.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
		}

		// Protected auth routes
		protectedAuthGroup := v1Group.Group("/auth")
		protectedAuthGroup.Use(middleware.AuthMiddleware(jwtService, l))
		{
			protectedAuthGroup.POST("/logout", authHandler.Logout)
			protectedAuthGroup.GET("/me", authHandler.Me)
		}

		// Protected API routes
		protectedGroup := v1Group.Group("")
		protectedGroup.Use(middleware.AuthMiddleware(jwtService, l))
		{
			v1.NewSpyCatsRoutes(protectedGroup, catHandlerService, l)
			v1.NewMissionsRoutes(protectedGroup, missionHandlerService, l)
			v1.NewTargetsRoutes(protectedGroup, targetHandlerService, l)
		}
	}
}
