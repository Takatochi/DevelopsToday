package auth

import (
	"context"
	"net/http"
	"time"

	"DevelopsToday/internal/controller/http/middleware"
	"DevelopsToday/internal/dto"
	"DevelopsToday/internal/models"
	"DevelopsToday/internal/repo"
	"DevelopsToday/internal/services"
	"DevelopsToday/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	userRepo   repo.UserRepository
	jwtService *services.JWTService
	logger     logger.Interface
}

func NewHandler(userRepo repo.UserRepository, jwtService *services.JWTService, logger logger.Interface) *Handler {
	return &Handler{
		userRepo:   userRepo,
		jwtService: jwtService,
		logger:     logger,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	ctx := context.Background()

	// Check if user already exists
	if _, err := h.userRepo.FindByUsername(ctx, req.Username); err == nil {
		_ = c.Error(middleware.ErrUserExists)
		return
	}

	if _, err := h.userRepo.FindByEmail(ctx, req.Email); err == nil {
		_ = c.Error(middleware.ErrEmailExists)
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "user"
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := h.userRepo.Create(ctx, user); err != nil {
		h.logger.Error("Failed to create user: %v", err)
		_ = c.Error(err)
		return
	}

	// Generate tokens
	tokens, err := h.jwtService.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		h.logger.Error("Failed to generate tokens: %v", err)
		_ = c.Error(middleware.ErrGenerateToken)
		return
	}

	response := dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	ctx := context.Background()

	// Find user
	user, err := h.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			_ = c.Error(middleware.ErrInvalidCreds)
			return
		}
		h.logger.Error("Failed to find user: %v", err)
		_ = c.Error(err)
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		_ = c.Error(middleware.ErrInvalidCreds)
		return
	}

	// Generate tokens
	tokens, err := h.jwtService.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		h.logger.Error("Failed to generate tokens: %v", err)
		_ = c.Error(middleware.ErrGenerateToken)
		return
	}

	response := dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	c.JSON(http.StatusOK, response)
}

// Refresh godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh token"
// @Success 200 {object} services.TokenPair
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(middleware.NewAppError("BAD_REQUEST", err.Error(), http.StatusBadRequest))
		return
	}

	tokens, err := h.jwtService.RefreshToken(req.RefreshToken)
	if err != nil {
		_ = c.Error(middleware.ErrInvalidToken)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout godoc
// @Summary Logout user
// @Description Revoke user's refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		_ = c.Error(middleware.ErrUnauthorized)
		h.logger.Warn("Logout failed missing user_id in context")
		return
	}

	if err := h.jwtService.RevokeToken(userID.(uint)); err != nil {
		_ = c.Error(middleware.ErrFailedRefresh)
		h.logger.Error("Failed to revoke token: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Me godoc
// @Summary Get current user
// @Description Get current authenticated user information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		_ = c.Error(middleware.ErrUnauthorized)
		return
	}

	ctx := context.Background()
	user, err := h.userRepo.FindByID(ctx, userID.(uint))
	if err != nil {
		h.logger.Error("Failed to find user: %v", err)
		_ = c.Error(err)
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
