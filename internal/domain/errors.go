package domain

import (
	"errors"
	"fmt"
)

// Common domain errors
var (
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrValidationFailed  = errors.New("validation failed")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternal          = errors.New("internal server error")
	ErrDatabaseOperation = errors.New("database operation failed")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrTimeout           = errors.New("operation timeout")
)

// ValidationError represents validation errors
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// NotFoundError represents resource not found errors
type NotFoundError struct {
	Resource string
	ID       string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%s' not found", e.Resource, e.ID)
}

// BusinessError represents business logic errors
type BusinessError struct {
	Code    string
	Message string
}

func (e BusinessError) Error() string {
	return fmt.Sprintf("business error [%s]: %s", e.Code, e.Message)
}

// DatabaseError represents database-specific errors
type DatabaseError struct {
	Operation string
	Err       error
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

func (e DatabaseError) Unwrap() error {
	return e.Err
}

// IsNotFound checks if error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidationFailed)
}

// IsInvalidInput checks if error is an invalid input error
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsDatabaseError checks if error is a database error
func IsDatabaseError(err error) bool {
	return errors.Is(err, ErrDatabaseOperation)
}

// IsRateLimitError checks if error is a rate limit error
func IsRateLimitError(err error) bool {
	return errors.Is(err, ErrRateLimitExceeded)
}

// IsTimeoutError checks if error is a timeout error
func IsTimeoutError(err error) bool {
	return errors.Is(err, ErrTimeout)
}
