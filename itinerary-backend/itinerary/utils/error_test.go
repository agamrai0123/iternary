package utils

import (
	"net/http"
	"testing"
)

// TestNewAPIError verifies APIError creation
func TestNewAPIError(t *testing.T) {
	tests := []struct {
		name     string
		code     ErrorCode
		message  string
		details  string
		wantCode ErrorCode
	}{
		{
			name:     "create invalid input error",
			code:     ErrInvalidInput,
			message:  "Invalid input provided",
			details:  "Field validation failed",
			wantCode: ErrInvalidInput,
		},
		{
			name:     "create not found error",
			code:     ErrNotFound,
			message:  "Resource not found",
			details:  "Trip ID does not exist",
			wantCode: ErrNotFound,
		},
		{
			name:     "create unauthorized error",
			code:     ErrUnauthorized,
			message:  "Unauthorized access",
			details:  "Invalid token",
			wantCode: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAPIError(tt.code, tt.message, tt.details)

			if err.Code != tt.wantCode {
				t.Errorf("Expected code %v, got %v", tt.wantCode, err.Code)
			}

			if err.Message != tt.message {
				t.Errorf("Expected message %q, got %q", tt.message, err.Message)
			}

			if err.Details != tt.details {
				t.Errorf("Expected details %q, got %q", tt.details, err.Details)
			}

			if err.Error() == "" {
				t.Error("Error() string representation should not be empty")
			}
		})
	}
}

// TestAPIErrorString verifies error string formatting
func TestAPIErrorString(t *testing.T) {
	err := NewAPIError(ErrValidationError, "Invalid data", "Email format is wrong")

	errStr := err.Error()

	if errStr == "" {
		t.Error("Error string should not be empty")
	}

	// Verify format includes code and message
	if !contains(errStr, string(ErrValidationError)) {
		t.Errorf("Error string should contain error code: %s", errStr)
	}

	if !contains(errStr, "Invalid data") {
		t.Errorf("Error string should contain message: %s", errStr)
	}
}

// TestGetStatusCode verifies HTTP status code mapping
func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		name          string
		errorCode     ErrorCode
		expectedCode  int
	}{
		{
			name:         "invalid input maps to 400",
			errorCode:    ErrInvalidInput,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation error maps to 400",
			errorCode:    ErrValidationError,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "not found maps to 404",
			errorCode:    ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "unauthorized maps to 401",
			errorCode:    ErrUnauthorized,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "forbidden maps to 403",
			errorCode:    ErrForbidden,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "conflict maps to 409",
			errorCode:    ErrConflict,
			expectedCode: http.StatusConflict,
		},
		{
			name:         "database error maps to 500",
			errorCode:    ErrDatabaseError,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "internal server error maps to 500",
			errorCode:    ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "file upload error maps to 400",
			errorCode:    ErrFileUploadError,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := GetStatusCode(tt.errorCode)

			if statusCode != tt.expectedCode {
				t.Errorf("GetStatusCode(%v) = %d, expected %d", tt.errorCode, statusCode, tt.expectedCode)
			}
		})
	}
}

// TestErrorCodes verifies all error codes are defined
func TestErrorCodes(t *testing.T) {
	codes := []ErrorCode{
		ErrInvalidInput,
		ErrNotFound,
		ErrUnauthorized,
		ErrForbidden,
		ErrConflict,
		ErrInternalServer,
		ErrDatabaseError,
		ErrValidationError,
		ErrFileUploadError,
	}

	for _, code := range codes {
		if code == "" {
			t.Error("Error code should not be empty")
		}

		// Verify each code maps to a valid HTTP status
		statusCode := GetStatusCode(code)
		if statusCode < 400 || statusCode >= 600 {
			t.Errorf("Invalid HTTP status code %d for error code %v", statusCode, code)
		}
	}
}

// TestAPIErrorWithTraceID verifies trace ID handling
func TestAPIErrorWithTraceID(t *testing.T) {
	err := NewAPIError(ErrDatabaseError, "Database connection failed", "Connection timeout")
	traceID := "trace-12345"
	err.TraceID = traceID

	if err.TraceID != traceID {
		t.Errorf("Expected trace ID %q, got %q", traceID, err.TraceID)
	}
}

// TestAPIErrorWithStatusCode verifies status code assignment
func TestAPIErrorWithStatusCode(t *testing.T) {
	err := NewAPIError(ErrNotFound, "Not found", "")
	err.StatusCode = http.StatusNotFound

	if err.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, err.StatusCode)
	}
}

// TestErrorCodeComparison verifies error codes can be compared
func TestErrorCodeComparison(t *testing.T) {
	err1 := NewAPIError(ErrValidationError, "Validation failed", "")
	err2 := NewAPIError(ErrValidationError, "Validation failed", "")
	err3 := NewAPIError(ErrNotFound, "Not found", "")

	if err1.Code != err2.Code {
		t.Error("Same error codes should be equal")
	}

	if err1.Code == err3.Code {
		t.Error("Different error codes should not be equal")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
