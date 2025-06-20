package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/Takatochi/DevelopsToday/internal/dto"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			var statusCode int
			var message string

			switch {
			case errors.Is(err.Err, ErrUnauthorized):
				statusCode = http.StatusUnauthorized
				message = "Unauthorized"
			case errors.Is(err.Err, ErrForbidden):
				statusCode = http.StatusForbidden
				message = "Forbidden"
			case errors.Is(err.Err, ErrNotFound):
				statusCode = http.StatusNotFound
				message = "Not found"
			case errors.Is(err.Err, ErrBadRequest):
				statusCode = http.StatusBadRequest
				message = err.Error()
			case errors.Is(err.Err, ErrValidation):
				statusCode = http.StatusBadRequest
				message = "Validation failed"
			case errors.Is(err.Err, ErrConflict):
				statusCode = http.StatusConflict
				message = "Conflict"
			default:
				statusCode = http.StatusInternalServerError
				message = "Internal server error"
			}

			c.JSON(statusCode, dto.ErrorResponse{
				Error:   message,
				Message: err.Error(),
			})
			c.Abort()
		}
	}
}

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrValidation   = errors.New("validation error")
	ErrConflict     = errors.New("conflict")
)
