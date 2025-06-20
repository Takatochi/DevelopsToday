package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// MemoryCacheItem represents a cached item with expiration
type MemoryCacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

// IsExpired checks if the cache item has expired
func (item *MemoryCacheItem) IsExpired() bool {
	return time.Now().After(item.ExpiresAt)
}

// MemoryCacheService implements CacheService using in-memory storage
type MemoryCacheService struct {
	data   map[string]*MemoryCacheItem
	mutex  sync.RWMutex
	ticker *time.Ticker
	done   chan bool
}

// NewMemoryCacheService creates a new in-memory cache service
func NewMemoryCacheService() CacheService {
	cache := &MemoryCacheService{
		data:   make(map[string]*MemoryCacheItem),
		ticker: time.NewTicker(time.Minute), // Clean up expired items every minute
		done:   make(chan bool),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// cleanup removes expired items from the cache
func (m *MemoryCacheService) cleanup() {
	for {
		select {
		case <-m.ticker.C:
			m.mutex.Lock()
			for key, item := range m.data {
				if item.IsExpired() {
					delete(m.data, key)
				}
			}
			m.mutex.Unlock()
		case <-m.done:
			return
		}
	}
}

// Set stores a key-value pair with optional TTL
func (m *MemoryCacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	expiresAt := time.Now().Add(ttl)
	if ttl <= 0 {
		expiresAt = time.Now().Add(time.Hour * 24 * 365) // 1 year for "no expiration"
	}

	m.data[key] = &MemoryCacheItem{
		Value:     value,
		ExpiresAt: expiresAt,
	}

	return nil
}

// Get retrieves a value by key
func (m *MemoryCacheService) Get(ctx context.Context, key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}

	if item.IsExpired() {
		// Remove expired item
		m.mutex.RUnlock()
		m.mutex.Lock()
		delete(m.data, key)
		m.mutex.Unlock()
		m.mutex.RLock()
		return "", fmt.Errorf("key expired")
	}

	// Convert value to string
	switch v := item.Value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		// Try to marshal as JSON
		data, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("failed to convert value to string: %w", err)
		}
		return string(data), nil
	}
}

// Delete removes a key from cache
func (m *MemoryCacheService) Delete(ctx context.Context, key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.data, key)
	return nil
}

// Exists checks if a key exists in cache
func (m *MemoryCacheService) Exists(ctx context.Context, key string) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	item, exists := m.data[key]
	if !exists {
		return false, nil
	}

	if item.IsExpired() {
		// Remove expired item
		m.mutex.RUnlock()
		m.mutex.Lock()
		delete(m.data, key)
		m.mutex.Unlock()
		m.mutex.RLock()
		return false, nil
	}

	return true, nil
}

// SetJSON stores a JSON-serializable object
func (m *MemoryCacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return m.Set(ctx, key, string(data), ttl)
}

// GetJSON retrieves and unmarshals a JSON object
func (m *MemoryCacheService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := m.Get(ctx, key)
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
func (m *MemoryCacheService) Ping(ctx context.Context) error {
	// Memory cache is always available
	return nil
}

// Close closes the cache connection
func (m *MemoryCacheService) Close() error {
	m.ticker.Stop()
	close(m.done)

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Clear all data
	m.data = make(map[string]*MemoryCacheItem)

	return nil
}

// Size returns the number of items in the cache
func (m *MemoryCacheService) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.data)
}

// Clear removes all items from the cache
func (m *MemoryCacheService) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = make(map[string]*MemoryCacheItem)
}
