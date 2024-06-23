package utils

import (
	"fmt"
)

// HTTPError represents a custom HTTP error.
type HTTPError struct {
	Message string
	Code    int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// NewUnauthorizedError creates a new UnauthorizedError (401).
func NewUnauthorizedError(message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    401,
	}
}

// NewBadRequestError creates a new BadRequestError (400).
func NewBadRequestError(message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    400,
	}
}

// NewConflictError creates a new ConflictError (409).
func NewConflictError(message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    409,
	}
}

// NewInternalServerError creates a new InternalServerError (500).
func NewInternalServerError(message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    500,
	}
}

// NewUnauthenticatedError creates a new UnauthenticatedError (401).
func NewUnauthenticatedError(message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    401,
	}
}

// NewNotFoundError creates a new NotFoundError (404).
func NewNotFoundError(message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    404,
	}
}

// NewValidationError creates a new ValidationError (400).
func NewValidationError(message string) *HTTPError {
	return &HTTPError{
		Message: message,
		Code:    400,
	}
}
