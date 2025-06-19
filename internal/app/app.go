package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"DevelopsToday/config"
	"DevelopsToday/internal/controller/http"
	"DevelopsToday/internal/repo"
	"DevelopsToday/internal/repo/postgres"
	"DevelopsToday/internal/services"
	"DevelopsToday/pkg/gorm"
	"DevelopsToday/pkg/logger"
	"DevelopsToday/pkg/server"

	"github.com/gin-gonic/gin"
	gormlog "gorm.io/gorm/logger"
)

// App represents the application
type App struct {
	Handler *gin.Engine
	Server  *server.Server
	Logger  *logger.Logger
}

// New creates a new application instance for testing
func New(cfg *config.Config) (*App, error) {
	l := logger.New(cfg.Log.Level)

	// Database
	db, err := gorm.Connect(cfg.PG.URL, gorm.WithLoggerLevel(gormlog.Silent))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = gorm.AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Repository
	store := postgres.NewRepository(db)

	// Cache Service
	cacheFactory := services.NewCacheFactory()
	cacheService, err := cacheFactory.CreateCacheService(services.CacheType(cfg.Cache.Type), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache service: %w", err)
	}

	// JWT Service
	jwtService := services.NewJWTService(cfg, cacheService)

	httpServer := server.New(
		server.Port(cfg.HTTP.Port),
	)

	http.NewV1Controller(httpServer.Engine, store, cfg, l, jwtService)

	return &App{
		Handler: httpServer.Engine,
		Server:  httpServer,
		Logger:  l,
	}, nil
}

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	l.Info("Starting application...")
	l.Info("Database URL: %s", cfg.PG.URL)

	// Database
	db, err := gorm.Connect(cfg.PG.URL, gorm.WithLoggerLevel(gormlog.Info))
	if err != nil {
		l.Error("Failed to connect to database: %v", err)
		panic(err)
	}
	l.Info("Database connected successfully")

	if err = gorm.AutoMigrate(db); err != nil {
		l.Error("Failed to migrate database: %v", err)
		panic(err)
	}
	l.Info("Database migration completed")

	// Seed database (ignore errors if data already exists)
	if err = repo.Seed(db); err != nil {
		l.Info("Database seeding skipped (data may already exist): %v", err)
	} else {
		l.Info("Database seeded successfully")
	}

	// Repository
	store := postgres.NewRepository(db)
	l.Info("Repository initialized")

	// Cache Service
	cacheFactory := services.NewCacheFactory()
	cacheService, err := cacheFactory.CreateCacheService(services.CacheType(cfg.Cache.Type), cfg)
	if err != nil {
		l.Error("Failed to create cache service: %v", err)
		panic(err)
	}
	l.Info("Cache service created and connected")

	// JWT Service
	jwtService := services.NewJWTService(cfg, cacheService)
	l.Info("JWT service initialized")

	httpServer := server.New(
		server.Port(cfg.HTTP.Port),
	)
	l.Info("HTTP server created on port: %s", cfg.HTTP.Port)

	http.NewV1Controller(httpServer.Engine, store, cfg, l, jwtService)
	l.Info("Controllers initialized")

	httpServer.Start()
	l.Info("HTTP server started successfully")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case serverErr := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", serverErr))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
