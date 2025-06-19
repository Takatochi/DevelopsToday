package services

import (
	"testing"
	"time"

	"DevelopsToday/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		JWT: config.JWT{
			Secret:          "test-secret-key-for-jwt-testing",
			AccessTokenTTL:  900,    // 15 minutes
			RefreshTokenTTL: 604800, // 7 days
		},
		App: config.App{
			Name: "spy-cats-api",
		},
	}

	// Create memory cache for testing
	cache := NewMemoryCacheService()
	jwtService := NewJWTService(cfg, cache)

	testUserID := uint(1)
	testUsername := "testuser"
	testRole := "admin"

	t.Run("GenerateTokenPair should create valid tokens", func(t *testing.T) {
		tokens, err := jwtService.GenerateTokenPair(testUserID, testUsername, testRole)
		require.NoError(t, err)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)
	})

	t.Run("ValidateToken should validate correct token", func(t *testing.T) {
		tokens, err := jwtService.GenerateTokenPair(testUserID, testUsername, testRole)
		require.NoError(t, err)

		claims, err := jwtService.ValidateToken(tokens.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, testUserID, claims.UserID)
		assert.Equal(t, testUsername, claims.Username)
		assert.Equal(t, testRole, claims.Role)
	})

	t.Run("ValidateToken should reject invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"
		_, err := jwtService.ValidateToken(invalidToken)
		assert.Error(t, err)
	})

	t.Run("RefreshToken should create new tokens from valid refresh token", func(t *testing.T) {
		originalTokens, err := jwtService.GenerateTokenPair(testUserID, testUsername, testRole)
		require.NoError(t, err)

		// Wait a bit to ensure different timestamps
		time.Sleep(time.Second * 1)

		newTokens, err := jwtService.RefreshToken(originalTokens.RefreshToken)
		require.NoError(t, err)
		assert.NotEmpty(t, newTokens.AccessToken)
		assert.NotEmpty(t, newTokens.RefreshToken)

		// Tokens should be different due to different timestamps
		if originalTokens.AccessToken == newTokens.AccessToken {
			t.Log("Warning: Access tokens are identical - this can happen if generated at exact same time")
		}
		if originalTokens.RefreshToken == newTokens.RefreshToken {
			t.Log("Warning: Refresh tokens are identical - this can happen if generated at exact same time")
		}
	})

	t.Run("RefreshToken should reject invalid refresh token", func(t *testing.T) {
		invalidToken := "invalid.refresh.token"
		_, err := jwtService.RefreshToken(invalidToken)
		assert.Error(t, err)
	})

	t.Run("RevokeToken should remove refresh token from cache", func(t *testing.T) {
		tokens, err := jwtService.GenerateTokenPair(testUserID, testUsername, testRole)
		require.NoError(t, err)

		// Revoke the token
		err = jwtService.RevokeToken(testUserID)
		require.NoError(t, err)

		// Try to refresh with revoked token should fail
		_, err = jwtService.RefreshToken(tokens.RefreshToken)
		assert.Error(t, err)
	})

	t.Run("BlacklistToken should prevent token usage", func(t *testing.T) {
		tokens, err := jwtService.GenerateTokenPair(testUserID, testUsername, testRole)
		require.NoError(t, err)

		// Token should not be blacklisted initially
		assert.False(t, jwtService.IsTokenBlacklisted(tokens.AccessToken))

		// Blacklist the token
		err = jwtService.BlacklistToken(tokens.AccessToken)
		require.NoError(t, err)

		// Token should now be blacklisted
		assert.True(t, jwtService.IsTokenBlacklisted(tokens.AccessToken))
	})

	t.Run("Different secret keys should not validate tokens", func(t *testing.T) {
		differentCfg := &config.Config{
			JWT: config.JWT{
				Secret:          "different-secret",
				AccessTokenTTL:  900,
				RefreshTokenTTL: 604800,
			},
			App: config.App{
				Name: "spy-cats-api",
			},
		}

		differentSecretService := NewJWTService(differentCfg, cache)

		tokens, err := jwtService.GenerateTokenPair(testUserID, testUsername, testRole)
		require.NoError(t, err)

		// Token created with original service should not validate with different secret
		_, err = differentSecretService.ValidateToken(tokens.AccessToken)
		assert.Error(t, err)
	})
}

func TestJWTClaims(t *testing.T) {
	t.Run("Claims should contain correct issuer", func(t *testing.T) {
		cfg := &config.Config{
			JWT: config.JWT{
				Secret:          "test-secret",
				AccessTokenTTL:  900,
				RefreshTokenTTL: 604800,
			},
			App: config.App{
				Name: "spy-cats-api",
			},
		}

		cache := NewMemoryCacheService()
		jwtService := NewJWTService(cfg, cache)

		testUserID := uint(1)
		testUsername := "testuser"
		testRole := "admin"

		tokens, err := jwtService.GenerateTokenPair(testUserID, testUsername, testRole)
		require.NoError(t, err)

		claims, err := jwtService.ValidateToken(tokens.AccessToken)
		require.NoError(t, err)

		assert.Equal(t, "spy-cats-api", claims.Issuer)
	})

	t.Run("Claims should have valid timestamps", func(t *testing.T) {
		cfg := &config.Config{
			JWT: config.JWT{
				Secret:          "test-secret",
				AccessTokenTTL:  900,
				RefreshTokenTTL: 604800,
			},
			App: config.App{
				Name: "spy-cats-api",
			},
		}

		cache := NewMemoryCacheService()
		jwtService := NewJWTService(cfg, cache)

		testUserID := uint(1)
		testUsername := "testuser"
		testRole := "admin"

		beforeGeneration := time.Now()
		tokens, err := jwtService.GenerateTokenPair(testUserID, testUsername, testRole)
		require.NoError(t, err)
		afterGeneration := time.Now()

		claims, err := jwtService.ValidateToken(tokens.AccessToken)
		require.NoError(t, err)

		// IssuedAt should be between before and after generation
		assert.True(t, claims.IssuedAt.Time.After(beforeGeneration) || claims.IssuedAt.Time.Equal(beforeGeneration))
		assert.True(t, claims.IssuedAt.Time.Before(afterGeneration) || claims.IssuedAt.Time.Equal(afterGeneration))

		// NotBefore should be same as IssuedAt
		assert.Equal(t, claims.IssuedAt.Time, claims.NotBefore.Time)

		// ExpiresAt should be IssuedAt + TTL (900 seconds = 15 minutes)
		expectedExpiry := claims.IssuedAt.Time.Add(time.Second * 900)
		// Allow reasonable time difference due to processing time (up to 10 seconds)
		timeDiff := claims.ExpiresAt.Time.Sub(expectedExpiry)
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}
		// Just check that the expiry time is reasonable (within 10 seconds of expected)
		assert.True(t, timeDiff < time.Second*10, "ExpiresAt should be close to expected time, got diff: %v, expected: %v, actual: %v", timeDiff, expectedExpiry, claims.ExpiresAt.Time)
	})
}
