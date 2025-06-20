package dto

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	// Error message
	// @example "Invalid request parameters"
	Error string `json:"error" example:"Invalid request parameters"`
	
	// Error code (optional)
	// @example "VALIDATION_ERROR"
	Code string `json:"code,omitempty" example:"VALIDATION_ERROR"`
	
	// Additional error details (optional)
	// @example {"field": "username", "message": "Username is required"}
	Details interface{} `json:"details,omitempty"`
}

// SuccessResponse represents a success response
// @Description Success response structure
type SuccessResponse struct {
	// Success message
	// @example "Operation completed successfully"
	Message string `json:"message" example:"Operation completed successfully"`
	
	// Response data (optional)
	Data interface{} `json:"data,omitempty"`
}

// PaginationMeta represents pagination metadata
// @Description Pagination information
type PaginationMeta struct {
	// Current page number
	// @example 1
	Page int `json:"page" example:"1"`
	
	// Number of items per page
	// @example 10
	Limit int `json:"limit" example:"10"`
	
	// Total number of items
	// @example 100
	Total int64 `json:"total" example:"100"`
	
	// Total number of pages
	// @example 10
	TotalPages int `json:"total_pages" example:"10"`
}

// PaginatedResponse represents a paginated response
// @Description Paginated response structure
type PaginatedResponse struct {
	// Response data
	Data interface{} `json:"data"`
	
	// Pagination metadata
	Meta PaginationMeta `json:"meta"`
}

// HealthResponse represents health check response
// @Description Health check response
type HealthResponse struct {
	// Service status
	// @example "OK"
	Status string `json:"status" example:"OK"`
	
	// Timestamp of the health check
	// @example "2023-12-01T10:00:00Z"
	Timestamp string `json:"timestamp,omitempty" example:"2023-12-01T10:00:00Z"`
	
	// Service version (optional)
	// @example "1.0.0"
	Version string `json:"version,omitempty" example:"1.0.0"`
}
