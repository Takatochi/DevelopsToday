package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		JWT     JWT
		App     App
		PG      PG
		Log     Log
		HTTP    HTTP
		Swagger Swagger
		Cache   Cache
		Redis   Redis
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	PG struct {
		URL     string `env:"PG_URL,required"`
		PoolMax int    `env:"PG_POOL_MAX,required"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
	}

	JWT struct {
		Secret           string `env:"JWT_SECRET,required"`
		SigningAlgorithm string `env:"JWT_SIGNING_ALGORITHM" envDefault:"HS256"`
		AccessTokenTTL   int    `env:"JWT_ACCESS_TOKEN_TTL" envDefault:"900"`
		RefreshTokenTTL  int    `env:"JWT_REFRESH_TOKEN_TTL" envDefault:"604800"`
	}

	Cache struct {
		Type string `env:"CACHE_TYPE" envDefault:"redis"`
	}

	Redis struct {
		URL      string `env:"REDIS_URL" envDefault:"redis://localhost:6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
