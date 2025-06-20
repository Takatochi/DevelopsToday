package middleware

import (
	"fmt"
	"net/http"

	"DevelopsToday/internal/dto"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func GlobalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appErr *AppError
			var statusCode int
			var errorCode string
			var message string

			switch e := err.(type) {
			case *AppError:
				appErr = e
				statusCode = appErr.Status
				errorCode = appErr.Code
				message = appErr.Message
			case *ValidationError:
				statusCode = http.StatusBadRequest
				errorCode = "VALIDATION_ERROR"
				message = e.Error()
			case *AuthError:
				statusCode = e.Status
				errorCode = e.Code
				message = e.Message
			case *BusinessError:
				statusCode = e.Status
				errorCode = e.Code
				message = e.Message
			default:
				statusCode = http.StatusInternalServerError
				errorCode = "INTERNAL_ERROR"
				message = "Internal server error"
			}

			c.JSON(statusCode, dto.ErrorResponse{
				Error: message,
				Code:  errorCode,
			})
			c.Abort()
		}
	}
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

type AuthError struct {
	Code    string
	Message string
	Status  int
}

func (e *AuthError) Error() string {
	return e.Message
}

func NewAuthError(code, message string, status int) *AuthError {
	return &AuthError{Code: code, Message: message, Status: status}
}

type BusinessError struct {
	Code    string
	Message string
	Status  int
}

func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(code, message string, status int) *BusinessError {
	return &BusinessError{Code: code, Message: message, Status: status}
}

// Predefined errors
var (
	ErrUnauthorized = NewAuthError("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized)
	ErrInvalidToken = NewAuthError("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized)
	ErrForbidden    = NewAuthError("FORBIDDEN", "Access denied", http.StatusForbidden)
	ErrInvalidCreds = NewAuthError("INVALID_CREDENTIALS", "Invalid username or password", http.StatusUnauthorized)

	ErrNotFound        = NewAppError("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrUserNotFound    = NewAppError("USER_NOT_FOUND", "User not found", http.StatusNotFound)
	ErrCatNotFound     = NewAppError("CAT_NOT_FOUND", "Cat not found", http.StatusNotFound)
	ErrMissionNotFound = NewAppError("MISSION_NOT_FOUND", "Mission not found", http.StatusNotFound)
	ErrTargetNotFound  = NewAppError("TARGET_NOT_FOUND", "Target not found", http.StatusNotFound)

	ErrConflict    = NewAppError("CONFLICT", "Resource already exists", http.StatusConflict)
	ErrUserExists  = NewAppError("USER_EXISTS", "User already exists", http.StatusConflict)
	ErrEmailExists = NewAppError("EMAIL_EXISTS", "Email already registered", http.StatusConflict)

	ErrBadRequest   = NewAppError("BAD_REQUEST", "Invalid request", http.StatusBadRequest)
	ErrInvalidInput = NewAppError("INVALID_INPUT", "Invalid input data", http.StatusBadRequest)
	ErrMissingField = NewAppError("MISSING_FIELD", "Required field is missing", http.StatusBadRequest)

	ErrCatBusy         = NewBusinessError("CAT_BUSY", "Cat is already assigned to another mission", http.StatusConflict)
	ErrMissionComplete = NewBusinessError("MISSION_COMPLETE", "Mission is already completed", http.StatusBadRequest)
	ErrTargetComplete  = NewBusinessError("TARGET_COMPLETE", "Target is already completed", http.StatusBadRequest)
	ErrInvalidBreed    = NewBusinessError("INVALID_BREED", "Invalid cat breed", http.StatusBadRequest)
)
