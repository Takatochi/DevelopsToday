package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"DevelopsToday/internal/dto"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized  = NewAuthError("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized)
	ErrInvalidToken  = NewAuthError("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized)
	ErrForbidden     = NewAuthError("FORBIDDEN", "Access denied", http.StatusForbidden)
	ErrInvalidCreds  = NewAuthError("INVALID_CREDENTIALS", "Invalid username or password", http.StatusUnauthorized)
	ErrFailedRefresh = NewAuthError("FAILED_REFRESH", "Failed to refresh token", http.StatusUnauthorized)
	ErrRevokedToken  = NewAuthError("REVOKED_TOKEN", "Token has been revoked", http.StatusUnauthorized)
	ErrGenerateToken = NewAuthError("FAILED_GENERATE_TOKEN", "Failed to generate token", http.StatusInternalServerError)

	ErrNotFound        = NewAppError("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrUserNotFound    = NewAppError("USER_NOT_FOUND", "User not found", http.StatusNotFound)
	ErrCatNotFound     = NewAppError("CAT_NOT_FOUND", "Cat not found", http.StatusNotFound)
	ErrMissionNotFound = NewAppError("MISSION_NOT_FOUND", "Mission not found", http.StatusNotFound)
	ErrTargetNotFound  = NewAppError("TARGET_NOT_FOUND", "Target not found", http.StatusNotFound)

	ErrConflict    = NewAppError("CONFLICT", "Resource already exists", http.StatusConflict)
	ErrUserExists  = NewAppError("USER_EXISTS", "User already exists", http.StatusConflict)
	ErrEmailExists = NewAppError("EMAIL_EXISTS", "Email already registered", http.StatusConflict)

	ErrBadRequest    = NewAppError("BAD_REQUEST", "Invalid request", http.StatusBadRequest)
	ErrInvalidInput  = NewAppError("INVALID_INPUT", "Invalid input data", http.StatusBadRequest)
	ErrMissingField  = NewAppError("MISSING_FIELD", "Required field is missing", http.StatusBadRequest)
	ErrInternalError = NewAppError("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)

	ErrCatBusy         = NewBusinessError("CAT_BUSY", "Cat is already assigned to another mission", http.StatusConflict)
	ErrMissionComplete = NewBusinessError("MISSION_COMPLETE", "Mission is already completed", http.StatusBadRequest)
	ErrTargetComplete  = NewBusinessError("TARGET_COMPLETE", "Target is already completed", http.StatusBadRequest)
	ErrInvalidBreed    = NewBusinessError("INVALID_BREED", "Invalid cat breed", http.StatusBadRequest)

	StatusCodeValidation = "VALIDATION_ERROR"
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
			appErr := errorHandler(err)
			c.JSON(appErr.Status, dto.ErrorResponse{
				Error: appErr.Message,
				Code:  appErr.Code,
			})
			c.Abort()
		}
	}
}
func errorHandler(err error) *AppError {
	var appErr *AppError
	var valErr *ValidationError
	var authErr *AuthError
	var bizErr *BusinessError

	switch {
	case errors.As(err, &appErr):
		return &AppError{
			Code:    appErr.Code,
			Message: appErr.Message,
			Status:  appErr.Status,
		}
	case errors.As(err, &valErr):
		return &AppError{
			Code:    StatusCodeValidation,
			Message: valErr.Error(),
			Status:  http.StatusBadRequest,
		}
	case errors.As(err, &authErr):
		return &AppError{
			Code:    authErr.Code,
			Message: authErr.Message,
			Status:  authErr.Status,
		}
	case errors.As(err, &bizErr):
		return &AppError{
			Code:    bizErr.Code,
			Message: bizErr.Message,
			Status:  bizErr.Status,
		}
	default:
		return &AppError{
			Code:    ErrInternalError.Code,
			Message: ErrInternalError.Message,
			Status:  ErrInternalError.Status,
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
