package validation

import (
	"testing"
)

func TestValidator_ValidateNews(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content string
		want    bool
	}{
		{
			name:    "valid news",
			title:   "Valid Title",
			content: "This is a valid content with more than 10 characters",
			want:    true,
		},
		{
			name:    "empty title",
			title:   "",
			content: "Valid content",
			want:    false,
		},
		{
			name:    "title too short",
			title:   "Hi",
			content: "Valid content",
			want:    false,
		},
		{
			name:    "title too long",
			title:   string(make([]byte, 201)),
			content: "Valid content",
			want:    false,
		},
		{
			name:    "empty content",
			title:   "Valid Title",
			content: "",
			want:    false,
		},
		{
			name:    "content too short",
			title:   "Valid Title",
			content: "Short",
			want:    false,
		},
		{
			name:    "whitespace only title",
			title:   "   ",
			content: "Valid content",
			want:    false,
		},
		{
			name:    "whitespace only content",
			title:   "Valid Title",
			content: "   ",
			want:    false,
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateNews(tt.title, tt.content)
			if result.IsValid != tt.want {
				t.Errorf("ValidateNews() = %v, want %v", result.IsValid, tt.want)
			}
		})
	}
}

func TestValidator_ValidatePagination(t *testing.T) {
	tests := []struct {
		name  string
		page  int
		limit int
		want  bool
	}{
		{
			name:  "valid pagination",
			page:  1,
			limit: 10,
			want:  true,
		},
		{
			name:  "page zero",
			page:  0,
			limit: 10,
			want:  false,
		},
		{
			name:  "page negative",
			page:  -1,
			limit: 10,
			want:  false,
		},
		{
			name:  "limit zero",
			page:  1,
			limit: 0,
			want:  false,
		},
		{
			name:  "limit negative",
			page:  1,
			limit: -1,
			want:  false,
		},
		{
			name:  "limit too high",
			page:  1,
			limit: 101,
			want:  false,
		},
		{
			name:  "limit at boundary",
			page:  1,
			limit: 100,
			want:  true,
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidatePagination(tt.page, tt.limit)
			if result.IsValid != tt.want {
				t.Errorf("ValidatePagination() = %v, want %v", result.IsValid, tt.want)
			}
		})
	}
}

func TestValidator_ValidateSearchQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{
			name:  "valid query",
			query: "search term",
			want:  true,
		},
		{
			name:  "empty query",
			query: "",
			want:  false,
		},
		{
			name:  "whitespace only",
			query: "   ",
			want:  false,
		},
		{
			name:  "single character",
			query: "a",
			want:  true,
		},
	}

	validator := NewValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateSearchQuery(tt.query)
			if result.IsValid != tt.want {
				t.Errorf("ValidateSearchQuery() = %v, want %v", result.IsValid, tt.want)
			}
		})
	}
}

func TestValidator_Error(t *testing.T) {
	validator := NewValidator()

	// Test with valid result
	result := ValidationResult{IsValid: true}
	err := validator.Error(result)
	if err != nil {
		t.Errorf("Error() should return nil for valid result, got %v", err)
	}

	// Test with invalid result
	result = ValidationResult{
		IsValid: false,
		Errors: []ValidationError{
			{Field: "title", Message: "cannot be empty"},
			{Field: "content", Message: "too short"},
		},
	}
	err = validator.Error(result)
	if err == nil {
		t.Error("Error() should return error for invalid result")
	}

	expected := "validation failed: title: cannot be empty; content: too short"
	if err.Error() != expected {
		t.Errorf("Error() = %v, want %v", err.Error(), expected)
	}
}
