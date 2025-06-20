package services

import (
	"testing"

	"DevelopsToday/config"

	"github.com/stretchr/testify/assert"
)

func TestCacheFactory(t *testing.T) {
	factory := NewCacheFactory()

	t.Run("CreateCacheService with redis type should fail without valid Redis", func(t *testing.T) {
		cfg := &config.Config{
			Redis: config.Redis{
				URL:      "redis://invalid-host:6379",
				Password: "",
				DB:       0,
			},
		}

		cache, err := factory.CreateCacheService(CacheTypeRedis, cfg)
		assert.Error(t, err)
		assert.Nil(t, cache)
	})

	t.Run("CreateCacheService with memcached type should return error (not implemented)", func(t *testing.T) {
		cfg := &config.Config{}

		cache, err := factory.CreateCacheService(CacheTypeMemcached, cfg)
		assert.Error(t, err)
		assert.Nil(t, cache)
		assert.Contains(t, err.Error(), "not yet implemented")
	})

	t.Run("CreateCacheService with invalid type should return error", func(t *testing.T) {
		cfg := &config.Config{}

		cache, err := factory.CreateCacheService("invalid", cfg)
		assert.Error(t, err)
		assert.Nil(t, cache)
		assert.Contains(t, err.Error(), "unsupported cache type")
	})

	t.Run("createMemcachedCache should return not implemented error", func(t *testing.T) {
		cfg := &config.Config{}

		cache, err := factory.createMemcachedCache(cfg)
		assert.Error(t, err)
		assert.Nil(t, cache)
		assert.Contains(t, err.Error(), "not yet implemented")
	})
}

func TestCacheFactoryRedisURLParsing(t *testing.T) {
	factory := NewCacheFactory()

	t.Run("Redis URL parsing should handle invalid host", func(t *testing.T) {
		cfg := &config.Config{
			Redis: config.Redis{
				URL:      "redis://invalid-nonexistent-host:6379/0",
				Password: "",
				DB:       0,
			},
		}

		// This should fail to connect to invalid host
		cache, err := factory.CreateCacheService(CacheTypeRedis, cfg)
		assert.Error(t, err) // Expected to fail due to invalid Redis server
		assert.Nil(t, cache)
		assert.Contains(t, err.Error(), "failed to connect to Redis")
	})

	t.Run("Redis URL parsing should handle invalid port", func(t *testing.T) {
		cfg := &config.Config{
			Redis: config.Redis{
				URL:      "redis://localhost:99999/0",
				Password: "",
				DB:       0,
			},
		}

		// This should fail to connect to invalid port
		cache, err := factory.CreateCacheService(CacheTypeRedis, cfg)
		assert.Error(t, err) // Expected to fail due to invalid port
		assert.Nil(t, cache)
		assert.Contains(t, err.Error(), "failed to connect to Redis")
	})
}
