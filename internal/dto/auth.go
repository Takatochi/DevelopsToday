package dto

// RegisterRequest represents the registration request
// @Description User registration request
type RegisterRequest struct {
	// Username for the new user account
	// @example "john_doe"
	Username string `json:"username" binding:"required,min=3,max=50" example:"john_doe"`
	
	// Email address for the new user account
	// @example "john.doe@example.com"
	Email string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	
	// Password for the new user account (minimum 6 characters)
	// @example "securepassword123"
	Password string `json:"password" binding:"required,min=6" example:"securepassword123"`
	
	// Role for the user (optional, defaults to "spy")
	// @example "spy"
	Role string `json:"role,omitempty" example:"spy"`
}

// LoginRequest represents the login request
// @Description User login request
type LoginRequest struct {
	// Username or email for authentication
	// @example "john_doe"
	Username string `json:"username" binding:"required" example:"john_doe"`
	
	// Password for authentication
	// @example "securepassword123"
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

// RefreshRequest represents the refresh token request
// @Description Refresh token request
type RefreshRequest struct {
	// Refresh token to obtain new access token
	// @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// AuthResponse represents the authentication response
// @Description Authentication response with user data and tokens
type AuthResponse struct {
	// User information
	User UserResponse `json:"user"`
	
	// JWT access token for API authentication
	// @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	
	// JWT refresh token for obtaining new access tokens
	// @example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserResponse represents user data in responses
// @Description User information in API responses
type UserResponse struct {
	// Unique user identifier
	// @example 1
	ID uint `json:"id" example:"1"`
	
	// Username of the user
	// @example "john_doe"
	Username string `json:"username" example:"john_doe"`
	
	// Email address of the user
	// @example "john.doe@example.com"
	Email string `json:"email" example:"john.doe@example.com"`
	
	// Role of the user
	// @example "spy"
	Role string `json:"role" example:"spy"`
	
	// Account creation timestamp
	// @example "2023-12-01T10:00:00Z"
	CreatedAt string `json:"created_at" example:"2023-12-01T10:00:00Z"`
	
	// Last update timestamp
	// @example "2023-12-01T10:00:00Z"
	UpdatedAt string `json:"updated_at" example:"2023-12-01T10:00:00Z"`
}
