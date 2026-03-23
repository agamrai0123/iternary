package itinerary

import (
	"fmt"
	"net/http"
)

// ErrorCode represents an error code for the API
type ErrorCode string

const (
	ErrInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrNotFound        ErrorCode = "NOT_FOUND"
	ErrUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrForbidden       ErrorCode = "FORBIDDEN"
	ErrConflict        ErrorCode = "CONFLICT"
	ErrInternalServer  ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrDatabaseError   ErrorCode = "DATABASE_ERROR"
	ErrValidationError ErrorCode = "VALIDATION_ERROR"
	ErrFileUploadError ErrorCode = "FILE_UPLOAD_ERROR"
)

// APIError represents a structured API error response
type APIError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	StatusCode int       `json:"status_code"`
	Details    string    `json:"details,omitempty"`
	TraceID    string    `json:"trace_id,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
}

// NewAPIError creates a new API error
func NewAPIError(code ErrorCode, message string, details string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// GetStatusCode returns the HTTP status code for the error
func GetStatusCode(code ErrorCode) int {
	switch code {
	case ErrInvalidInput, ErrValidationError:
		return http.StatusBadRequest
	case ErrNotFound:
		return http.StatusNotFound
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrConflict:
		return http.StatusConflict
	case ErrFileUploadError:
		return http.StatusBadRequest
	case ErrDatabaseError, ErrInternalServer:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// NewInvalidInputError creates a validation error
func NewInvalidInputError(field string, reason string) *APIError {
	err := NewAPIError(
		ErrInvalidInput,
		"Invalid input",
		fmt.Sprintf("Field '%s': %s", field, reason),
	)
	err.StatusCode = http.StatusBadRequest
	return err
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string, id string) *APIError {
	err := NewAPIError(
		ErrNotFound,
		fmt.Sprintf("%s not found", resource),
		fmt.Sprintf("%s with ID '%s' does not exist", resource, id),
	)
	err.StatusCode = http.StatusNotFound
	return err
}

// NewDatabaseError creates a database error
func NewDatabaseError(operation string, err error) *APIError {
	apiErr := NewAPIError(
		ErrDatabaseError,
		fmt.Sprintf("Database %s failed", operation),
		err.Error(),
	)
	apiErr.StatusCode = http.StatusInternalServerError
	return apiErr
}

// NewValidationError creates a validation error
func NewValidationError(message string, details string) *APIError {
	err := NewAPIError(
		ErrValidationError,
		message,
		details,
	)
	err.StatusCode = http.StatusBadRequest
	return err
}

// NewInternalServerError creates an internal server error
func NewInternalServerError(operation string, err error) *APIError {
	apiErr := NewAPIError(
		ErrInternalServer,
		fmt.Sprintf("%s failed", operation),
		err.Error(),
	)
	apiErr.StatusCode = http.StatusInternalServerError
	return apiErr
}

// NewAuthenticationError creates an authentication error
func NewAuthenticationError(message string) *APIError {
	return NewAPIError(ErrUnauthorized, "Authentication Failed", message)
}

// NewAuthorizationError creates an authorization error
func NewAuthorizationError(message string) *APIError {
	return NewAPIError(ErrUnauthorized, "Authorization Failed", message)
}

// ToJSON converts error to JSON response
func (e *APIError) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    e.Code,
			"message": e.Message,
			"details": e.Details,
		},
	}
}
