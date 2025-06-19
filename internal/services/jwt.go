package services

import (
	"context"
	"fmt"
	"time"

	"DevelopsToday/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	cfg    *config.Config
	cache  CacheService
	secret []byte
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewJWTService(cfg *config.Config, cacheService CacheService) *JWTService {
	return &JWTService{
		cfg:    cfg,
		cache:  cacheService,
		secret: []byte(cfg.JWT.Secret),
	}
}

// GenerateTokenPair generates access and refresh tokens
func (j *JWTService) GenerateTokenPair(userID uint, username, role string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := j.generateToken(userID, username, role, time.Duration(j.cfg.JWT.AccessTokenTTL)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := j.generateToken(userID, username, role, time.Duration(j.cfg.JWT.RefreshTokenTTL)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in cache
	ctx := context.Background()
	refreshKey := fmt.Sprintf("refresh_token:%d", userID)
	err = j.cache.Set(ctx, refreshKey, refreshToken, time.Duration(j.cfg.JWT.RefreshTokenTTL)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateToken creates a JWT token
func (j *JWTService) generateToken(userID uint, username, role string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.cfg.App.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshToken refreshes an access token using a refresh token
func (j *JWTService) RefreshToken(refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if refresh token exists in cache
	ctx := context.Background()
	refreshKey := fmt.Sprintf("refresh_token:%d", claims.UserID)
	storedToken, err := j.cache.Get(ctx, refreshKey)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found: %w", err)
	}

	if storedToken != refreshToken {
		return nil, fmt.Errorf("refresh token mismatch")
	}

	// Generate new token pair
	return j.GenerateTokenPair(claims.UserID, claims.Username, claims.Role)
}

// RevokeToken revokes a refresh token
func (j *JWTService) RevokeToken(userID uint) error {
	ctx := context.Background()
	refreshKey := fmt.Sprintf("refresh_token:%d", userID)
	return j.cache.Delete(ctx, refreshKey)
}

// BlacklistToken adds a token to the blacklist
func (j *JWTService) BlacklistToken(tokenString string) error {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return err
	}

	ctx := context.Background()
	blacklistKey := fmt.Sprintf("blacklist:%s", tokenString)
	ttl := time.Until(claims.ExpiresAt.Time)

	return j.cache.Set(ctx, blacklistKey, "1", ttl)
}

// IsTokenBlacklisted checks if a token is blacklisted
func (j *JWTService) IsTokenBlacklisted(tokenString string) bool {
	ctx := context.Background()
	blacklistKey := fmt.Sprintf("blacklist:%s", tokenString)
	exists, err := j.cache.Exists(ctx, blacklistKey)
	return err == nil && exists
}
