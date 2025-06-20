package services

import (
	"context"
	"time"
)

// CacheService defines the interface for caching operations
type CacheService interface {
	// Set stores a key-value pair with optional TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Get retrieves a value by key
	Get(ctx context.Context, key string) (string, error)

	// Delete removes a key from cache
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)

	// SetJSON stores a JSON-serializable object
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// GetJSON retrieves and unmarshals a JSON object
	GetJSON(ctx context.Context, key string, dest interface{}) error

	// Ping checks if the cache service is available
	Ping(ctx context.Context) error

	// Close closes the cache connection
	Close() error
}
