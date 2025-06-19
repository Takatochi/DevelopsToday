package app

import (
	"DevelopsToday/config"
	"DevelopsToday/internal/controller/http"
	"DevelopsToday/pkg/logger"
	"DevelopsToday/pkg/server"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)


func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	httpServer := server.New(
		server.Port(cfg.HTTP.Port),
	)
	http.NewV1Controller(httpServer.Engine, cfg, l)
	httpServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
