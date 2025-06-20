package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCacheService implements CacheService using Redis
type RedisCacheService struct {
	client *redis.Client
}

// NewRedisCacheService creates a new Redis cache service
func NewRedisCacheService(client *redis.Client) CacheService {
	return &RedisCacheService{
		client: client,
	}
}

// Set stores a key-value pair with optional TTL
func (r *RedisCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves a value by key
func (r *RedisCacheService) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete removes a key from cache
func (r *RedisCacheService) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache
func (r *RedisCacheService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

// SetJSON stores a JSON-serializable object
func (r *RedisCacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return r.Set(ctx, key, string(data), ttl)
}

// GetJSON retrieves and unmarshals a JSON object
func (r *RedisCacheService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Get(ctx, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// Ping checks if the cache service is available
func (r *RedisCacheService) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the cache connection
func (r *RedisCacheService) Close() error {
	return r.client.Close()
}
