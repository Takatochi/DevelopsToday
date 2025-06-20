package services

import (
	"context"
	"fmt"
	"strings"

	"DevelopsToday/config"

	"github.com/redis/go-redis/v9"
)

// CacheType represents the type of cache service
type CacheType string

const (
	CacheTypeRedis     CacheType = "redis"
	CacheTypeMemcached CacheType = "memcached"
	CacheTypeMemory    CacheType = "memory"
)

// CacheFactory creates cache services based on configuration
type CacheFactory struct{}

// NewCacheFactory creates a new cache factory
func NewCacheFactory() *CacheFactory {
	return &CacheFactory{}
}

// CreateCacheService creates a cache service based on the specified type
func (f *CacheFactory) CreateCacheService(cacheType CacheType, cfg *config.Config) (CacheService, error) {
	switch cacheType {
	case CacheTypeRedis:
		return f.createRedisCache(cfg)
	case CacheTypeMemcached:
		return f.createMemcachedCache(cfg)
	case CacheTypeMemory:
		return f.createMemoryCache(cfg)
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheType)
	}
}

// createRedisCache creates a Redis cache service
func (f *CacheFactory) createRedisCache(cfg *config.Config) (CacheService, error) {
	// Extract host and port from Redis URL
	redisAddr := cfg.Redis.URL
	if len(redisAddr) > 8 && redisAddr[:8] == "redis://" {
		redisAddr = redisAddr[8:] // Remove "redis://" prefix
	}

	// Remove database number if present
	if idx := strings.Index(redisAddr, "/"); idx != -1 {
		redisAddr = redisAddr[:idx]
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		Network:  "tcp4", // Force IPv4
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return NewRedisCacheService(client), nil
}

// createMemcachedCache creates a Memcached cache service
func (f *CacheFactory) createMemcachedCache(cfg *config.Config) (CacheService, error) {
	// TODO: Implement when memcached support is needed
	return nil, fmt.Errorf("memcached cache service not yet implemented")
}

// createMemoryCache creates an in-memory cache service
func (f *CacheFactory) createMemoryCache(_ *config.Config) (CacheService, error) {
	return NewMemoryCacheService(), nil
}
