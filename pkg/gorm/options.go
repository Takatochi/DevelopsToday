package gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormOption func(*gorm.Config)

func WithLoggerLevel(level logger.LogLevel) GormOption {
	return func(cfg *gorm.Config) {
		cfg.Logger = logger.Default.LogMode(level)
	}
}
