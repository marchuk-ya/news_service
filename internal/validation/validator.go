package validation

import (
	"fmt"
	"strings"
)

// Validator provides validation functionality
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationResult contains validation results
type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// ValidateNews validates news entity
func (v *Validator) ValidateNews(title, content string) ValidationResult {
	var errors []ValidationError

	title = strings.TrimSpace(title)
	if title == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title cannot be empty",
		})
	} else if len(title) < 3 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title must be at least 3 characters long",
		})
	} else if len(title) > 200 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "title cannot exceed 200 characters",
		})
	}

	content = strings.TrimSpace(content)
	if content == "" {
		errors = append(errors, ValidationError{
			Field:   "content",
			Message: "content cannot be empty",
		})
	} else if len(content) < 10 {
		errors = append(errors, ValidationError{
			Field:   "content",
			Message: "content must be at least 10 characters long",
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidatePagination validates pagination parameters
func (v *Validator) ValidatePagination(page, limit int) ValidationResult {
	var errors []ValidationError

	if page < 1 {
		errors = append(errors, ValidationError{
			Field:   "page",
			Message: "page number must be positive",
		})
	}

	if limit < 1 {
		errors = append(errors, ValidationError{
			Field:   "limit",
			Message: "limit must be positive",
		})
	} else if limit > 100 {
		errors = append(errors, ValidationError{
			Field:   "limit",
			Message: "limit cannot exceed 100",
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// ValidateSearchQuery validates search query
func (v *Validator) ValidateSearchQuery(query string) ValidationResult {
	query = strings.TrimSpace(query)
	if query == "" {
		return ValidationResult{
			IsValid: false,
			Errors: []ValidationError{
				{
					Field:   "query",
					Message: "search query cannot be empty",
				},
			},
		}
	}

	return ValidationResult{
		IsValid: true,
		Errors:  nil,
	}
}

// Error returns formatted validation error
func (v *Validator) Error(result ValidationResult) error {
	if result.IsValid {
		return nil
	}

	var messages []string
	for _, err := range result.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}

	return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
}
