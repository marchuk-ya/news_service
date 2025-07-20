package service

import (
	"context"
	"fmt"
	"strings"

	"news_service/internal/domain"
	"news_service/internal/tools"
	"news_service/internal/validation"
)

// newsUseCase implements domain.NewsUseCase
type newsUseCase struct {
	repo      domain.NewsRepository
	validator *validation.Validator
}

// NewNewsUseCase creates a new news use case
func NewNewsUseCase(repo domain.NewsRepository) domain.NewsUseCase {
	return &newsUseCase{
		repo:      repo,
		validator: validation.NewValidator(),
	}
}

// CreateNews creates a new news article with validation and business logic
func (uc *newsUseCase) CreateNews(ctx context.Context, title, content string) (*domain.News, error) {
	// Validate input
	result := uc.validator.ValidateNews(title, content)
	if !result.IsValid {
		return nil, fmt.Errorf("%w: %v", domain.ErrValidationFailed, uc.validator.Error(result))
	}

	// Sanitize content
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	// Create news object
	news := &domain.News{
		Title:     title,
		Content:   content,
		CreatedAt: tools.GetCurrentTime(),
		UpdatedAt: tools.GetCurrentTime(),
	}

	// Save to repository
	if err := uc.repo.Create(ctx, news); err != nil {
		return nil, fmt.Errorf("%w: failed to create news: %v", domain.ErrDatabaseOperation, err)
	}

	return news, nil
}

// GetNewsByID retrieves a news article by ID
func (uc *newsUseCase) GetNewsByID(ctx context.Context, id string) (*domain.News, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: news ID cannot be empty", domain.ErrInvalidInput)
	}

	news, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}

	return news, nil
}

// GetAllNews retrieves all news articles with pagination
func (uc *newsUseCase) GetAllNews(ctx context.Context, page, limit int) ([]*domain.News, int64, error) {
	// Validate pagination parameters
	result := uc.validator.ValidatePagination(page, limit)
	if !result.IsValid {
		// Apply defaults if validation fails
		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}
		if limit > 100 {
			limit = 100
		}
	}

	news, total, err := uc.repo.GetAll(ctx, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: failed to get news: %v", domain.ErrDatabaseOperation, err)
	}

	return news, total, nil
}

// UpdateNews updates an existing news article
func (uc *newsUseCase) UpdateNews(ctx context.Context, id, title, content string) (*domain.News, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: news ID cannot be empty", domain.ErrInvalidInput)
	}

	// Validate input
	result := uc.validator.ValidateNews(title, content)
	if !result.IsValid {
		return nil, fmt.Errorf("%w: %v", domain.ErrValidationFailed, uc.validator.Error(result))
	}

	// Check if news exists
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}

	// Sanitize content
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	// Update fields
	existing.Title = title
	existing.Content = content
	existing.UpdatedAt = tools.GetCurrentTime()

	// Save to repository
	if err := uc.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("%w: failed to update news: %v", domain.ErrDatabaseOperation, err)
	}

	return existing, nil
}

// DeleteNews deletes a news article
func (uc *newsUseCase) DeleteNews(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: news ID cannot be empty", domain.ErrInvalidInput)
	}

	// Check if news exists before deletion
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: failed to delete news: %v", domain.ErrDatabaseOperation, err)
	}

	return nil
}

// SearchNews searches for news articles
func (uc *newsUseCase) SearchNews(ctx context.Context, query string, page, limit int) ([]*domain.News, int64, error) {
	// Validate search query
	result := uc.validator.ValidateSearchQuery(query)
	if !result.IsValid {
		return nil, 0, fmt.Errorf("%w: %v", domain.ErrInvalidInput, uc.validator.Error(result))
	}

	// Validate pagination parameters
	paginationResult := uc.validator.ValidatePagination(page, limit)
	if !paginationResult.IsValid {
		// Apply defaults if validation fails
		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 10
		}
		if limit > 100 {
			limit = 100
		}
	}

	news, total, err := uc.repo.SearchNews(ctx, query, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: failed to search news: %v", domain.ErrDatabaseOperation, err)
	}

	return news, total, nil
}
