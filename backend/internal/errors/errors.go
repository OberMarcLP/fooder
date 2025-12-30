package errors

import (
	"encoding/json"
	"net/http"

	"github.com/nomdb/backend/internal/logger"
)

// ErrorResponse represents a structured API error response
type ErrorResponse struct {
	Error   string `json:"error"`           // User-friendly error message
	Code    string `json:"code"`            // Error code for programmatic handling
	Status  int    `json:"status"`          // HTTP status code
	Details string `json:"details,omitempty"` // Additional details (optional)
}

// Error codes for programmatic handling
const (
	CodeValidationError   = "VALIDATION_ERROR"
	CodeNotFound          = "NOT_FOUND"
	CodeDuplicate         = "DUPLICATE_ENTRY"
	CodeUnauthorized      = "UNAUTHORIZED"
	CodeForbidden         = "FORBIDDEN"
	CodeInternalError     = "INTERNAL_ERROR"
	CodeBadRequest        = "BAD_REQUEST"
	CodeConflict          = "CONFLICT"
	CodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
)

// Common error constructors

// BadRequest creates a 400 Bad Request error
func BadRequest(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:  message,
		Code:   CodeBadRequest,
		Status: http.StatusBadRequest,
	}
}

// ValidationError creates a 400 Validation Error
func ValidationError(message string, details string) *ErrorResponse {
	return &ErrorResponse{
		Error:   message,
		Code:    CodeValidationError,
		Status:  http.StatusBadRequest,
		Details: details,
	}
}

// NotFound creates a 404 Not Found error
func NotFound(resource string) *ErrorResponse {
	return &ErrorResponse{
		Error:  resource + " not found",
		Code:   CodeNotFound,
		Status: http.StatusNotFound,
	}
}

// Duplicate creates a 409 Conflict error for duplicate entries
func Duplicate(resource string) *ErrorResponse {
	return &ErrorResponse{
		Error:  resource + " already exists",
		Code:   CodeDuplicate,
		Status: http.StatusConflict,
	}
}

// Conflict creates a 409 Conflict error
func Conflict(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:  message,
		Code:   CodeConflict,
		Status: http.StatusConflict,
	}
}

// InternalError creates a 500 Internal Server Error
func InternalError(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:  "Internal server error",
		Code:   CodeInternalError,
		Status: http.StatusInternalServerError,
		Details: message, // Log the actual error but don't expose to client in production
	}
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:  message,
		Code:   CodeUnauthorized,
		Status: http.StatusUnauthorized,
	}
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:  message,
		Code:   CodeForbidden,
		Status: http.StatusForbidden,
	}
}

// RespondWithError writes a JSON error response to the client
func RespondWithError(w http.ResponseWriter, err *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)

	if encodeErr := json.NewEncoder(w).Encode(err); encodeErr != nil {
		logger.Error("Failed to encode error response: %v", encodeErr)
	}

	// Log the error with context
	if err.Status >= 500 {
		logger.Error("Error %d: %s - %s (details: %s)", err.Status, err.Code, err.Error, err.Details)
	} else if err.Status >= 400 {
		logger.Warn("Error %d: %s - %s", err.Status, err.Code, err.Error)
	}
}

// HandleDatabaseError converts common database errors to structured responses
func HandleDatabaseError(err error, resource string) *ErrorResponse {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for common PostgreSQL errors
	switch {
	case contains(errMsg, "no rows"):
		return NotFound(resource)
	case contains(errMsg, "duplicate key"), contains(errMsg, "unique constraint"):
		return Duplicate(resource)
	case contains(errMsg, "foreign key"):
		return Conflict("Cannot perform operation due to related data")
	default:
		return InternalError(errMsg)
	}
}

// contains is a helper to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && s[0:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr) ||
		containsMiddle(s, substr))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
