package app

import (
	"DevelopsToday/config"
	"DevelopsToday/internal/controller/http"
	"DevelopsToday/internal/repo"
	"DevelopsToday/internal/repo/postgres"
	"DevelopsToday/pkg/gorm"
	"DevelopsToday/pkg/logger"
	"DevelopsToday/pkg/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	gormlog "gorm.io/gorm/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	//Database
	db, err := gorm.Connect(cfg.PG.URL, gorm.WithLoggerLevel(gormlog.Info))
	if err != nil {
		l.Error(err)
		return
	}
	if err = gorm.AutoMigrate(db); err != nil {
		l.Error(err)
		return
	}
	if err = repo.Seed(db); err != nil {
		l.Error(err)
		return
	}

	//Repository
	store := postgres.NewRepository(db)
	httpServer := server.New(
		server.Port(cfg.HTTP.Port),
	)
	http.NewV1Controller(httpServer.Engine, store, cfg, l)
	httpServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
