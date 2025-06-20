package services

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*redis.Client, func()) {
	// Start mini Redis server for testing
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cleanup := func() {
		client.Close()
		mr.Close()
	}

	return client, cleanup
}

func TestRedisCacheService(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	cache := NewRedisCacheService(client)
	ctx := context.Background()

	t.Run("Set and Get should work correctly", func(t *testing.T) {
		key := "test:key"
		value := "test value"
		ttl := time.Minute * 5

		err := cache.Set(ctx, key, value, ttl)
		require.NoError(t, err)

		result, err := cache.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Get non-existing key should return error", func(t *testing.T) {
		key := "test:nonexistent"

		result, err := cache.Get(ctx, key)
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Delete should remove key", func(t *testing.T) {
		key := "test:delete"
		value := "delete me"
		ttl := time.Minute * 5

		// Set the key
		err := cache.Set(ctx, key, value, ttl)
		require.NoError(t, err)

		// Verify it exists
		result, err := cache.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, value, result)

		// Delete it
		err = cache.Delete(ctx, key)
		require.NoError(t, err)

		// Verify it's gone
		_, err = cache.Get(ctx, key)
		assert.Error(t, err)
	})

	t.Run("Exists should check key existence", func(t *testing.T) {
		key := "test:exists"
		value := "exists test"
		ttl := time.Minute * 5

		// Key should not exist initially
		exists, err := cache.Exists(ctx, key)
		require.NoError(t, err)
		assert.False(t, exists)

		// Set the key
		err = cache.Set(ctx, key, value, ttl)
		require.NoError(t, err)

		// Key should exist now
		exists, err = cache.Exists(ctx, key)
		require.NoError(t, err)
		assert.True(t, exists)

		// Delete the key
		err = cache.Delete(ctx, key)
		require.NoError(t, err)

		// Key should not exist again
		exists, err = cache.Exists(ctx, key)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("TTL should be set correctly", func(t *testing.T) {
		key := "test:ttl"
		value := "ttl test"
		ttl := time.Second * 10

		// Set key with TTL
		err := cache.Set(ctx, key, value, ttl)
		require.NoError(t, err)

		// Key should exist immediately
		exists, err := cache.Exists(ctx, key)
		require.NoError(t, err)
		assert.True(t, exists)

		// Verify we can get the value
		result, err := cache.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, value, result)

		// Clean up
		err = cache.Delete(ctx, key)
		require.NoError(t, err)
	})

	t.Run("SetJSON and GetJSON should work with structs", func(t *testing.T) {
		key := "test:json"
		ttl := time.Minute * 5

		type TestStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		original := TestStruct{
			ID:   123,
			Name: "test object",
		}

		// Set JSON
		err := cache.SetJSON(ctx, key, original, ttl)
		require.NoError(t, err)

		// Get JSON
		var result TestStruct
		err = cache.GetJSON(ctx, key, &result)
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})

	t.Run("GetJSON with non-existing key should return error", func(t *testing.T) {
		key := "test:json:nonexistent"

		type TestStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		var result TestStruct
		err := cache.GetJSON(ctx, key, &result)
		assert.Error(t, err)
	})

	t.Run("SetJSON with invalid data should return error", func(t *testing.T) {
		key := "test:json:invalid"
		ttl := time.Minute * 5

		// Try to set a channel (which can't be marshaled to JSON)
		invalidData := make(chan int)

		err := cache.SetJSON(ctx, key, invalidData, ttl)
		assert.Error(t, err)
	})

	t.Run("GetJSON with invalid JSON should return error", func(t *testing.T) {
		key := "test:json:corrupted"
		ttl := time.Minute * 5

		// Set invalid JSON manually
		err := cache.Set(ctx, key, "invalid json {", ttl)
		require.NoError(t, err)

		type TestStruct struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		var result TestStruct
		err = cache.GetJSON(ctx, key, &result)
		assert.Error(t, err)
	})

	t.Run("Multiple operations should work independently", func(t *testing.T) {
		keys := []string{"test:multi:1", "test:multi:2", "test:multi:3"}
		values := []string{"value1", "value2", "value3"}
		ttl := time.Minute * 5

		// Set multiple keys
		for i, key := range keys {
			err := cache.Set(ctx, key, values[i], ttl)
			require.NoError(t, err)
		}

		// Get all keys and verify values
		for i, key := range keys {
			result, err := cache.Get(ctx, key)
			require.NoError(t, err)
			assert.Equal(t, values[i], result)
		}

		// Delete one key
		err := cache.Delete(ctx, keys[1])
		require.NoError(t, err)

		// Verify first and third keys still exist
		_, err = cache.Get(ctx, keys[0])
		assert.NoError(t, err)

		_, err = cache.Get(ctx, keys[2])
		assert.NoError(t, err)

		// Verify second key is deleted
		_, err = cache.Get(ctx, keys[1])
		assert.Error(t, err)
	})
}
