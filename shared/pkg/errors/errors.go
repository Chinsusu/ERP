package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// Common error constructors
var (
	// BadRequest creates a 400 error
	BadRequest = func(message string) *AppError {
		return New("BAD_REQUEST", message, http.StatusBadRequest)
	}

	// Unauthorized creates a 401 error
	Unauthorized = func(message string) *AppError {
		return New("UNAUTHORIZED", message, http.StatusUnauthorized)
	}

	// Forbidden creates a 403 error
	Forbidden = func(message string) *AppError {
		return New("FORBIDDEN", message, http.StatusForbidden)
	}

	// NotFound creates a 404 error
	NotFound = func(resource string) *AppError {
		return New("NOT_FOUND", fmt.Sprintf("%s not found", resource), http.StatusNotFound)
	}

	// Conflict creates a 409 error
	Conflict = func(message string) *AppError {
		return New("CONFLICT", message, http.StatusConflict)
	}

	// Internal creates a 500 error
	Internal = func(err error) *AppError {
		return Wrap(err, "INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
	}

	// ValidationError creates a 422 error
	ValidationError = func(message string) *AppError {
		return New("VALIDATION_ERROR", message, http.StatusUnprocessableEntity)
	}
)

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

// Error implements the error interface
func (v *ValidationErrors) Error() string {
	return "validation failed"
}

// Add adds a validation error for a field
func (v *ValidationErrors) Add(field, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string][]string)
	}
	v.Errors[field] = append(v.Errors[field], message)
}

// HasErrors returns true if there are validation errors
func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// NewValidationErrors creates a new ValidationErrors
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make(map[string][]string),
	}
}
