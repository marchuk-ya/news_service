package service

import (
	"context"
	"fmt"
	"strings"

	"news_service/internal/domain"
	"news_service/internal/tools"
	"news_service/internal/validation"
)

type newsService struct {
	repo      domain.NewsRepository
	validator *validation.Validator
}

func NewNewsService(repo domain.NewsRepository) domain.NewsService {
	return &newsService{
		repo:      repo,
		validator: validation.NewValidator(),
	}
}

func (s *newsService) Create(ctx context.Context, news *domain.News) error {
	// Validate input
	result := s.validator.ValidateNews(news.Title, news.Content)
	if !result.IsValid {
		return fmt.Errorf("%w: %v", domain.ErrValidationFailed, s.validator.Error(result))
	}

	// Sanitize content
	news.Title = strings.TrimSpace(news.Title)
	news.Content = strings.TrimSpace(news.Content)

	// Set timestamps
	now := tools.GetCurrentTime()
	news.CreatedAt = now
	news.UpdatedAt = now

	return s.repo.Create(ctx, news)
}

func (s *newsService) GetByID(ctx context.Context, id string) (*domain.News, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: news ID cannot be empty", domain.ErrInvalidInput)
	}

	news, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}

	return news, nil
}

func (s *newsService) GetAll(ctx context.Context, page, limit int) ([]*domain.News, int64, error) {
	// Validate pagination parameters
	result := s.validator.ValidatePagination(page, limit)
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

	return s.repo.GetAll(ctx, page, limit)
}

func (s *newsService) Update(ctx context.Context, news *domain.News) error {
	// Validate input
	result := s.validator.ValidateNews(news.Title, news.Content)
	if !result.IsValid {
		return fmt.Errorf("%w: %v", domain.ErrValidationFailed, s.validator.Error(result))
	}

	// Check if news exists
	existing, err := s.repo.GetByID(ctx, news.ID.Hex())
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}

	// Sanitize content
	news.Title = strings.TrimSpace(news.Title)
	news.Content = strings.TrimSpace(news.Content)

	// Preserve original creation time
	news.CreatedAt = existing.CreatedAt
	news.UpdatedAt = tools.GetCurrentTime()

	return s.repo.Update(ctx, news)
}

func (s *newsService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: news ID cannot be empty", domain.ErrInvalidInput)
	}

	// Check if news exists before deletion
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrNotFound, err)
	}

	return s.repo.Delete(ctx, id)
}

func (s *newsService) SearchNews(ctx context.Context, query string, page, limit int) ([]*domain.News, int64, error) {
	// Validate search query
	result := s.validator.ValidateSearchQuery(query)
	if !result.IsValid {
		return nil, 0, fmt.Errorf("%w: %v", domain.ErrInvalidInput, s.validator.Error(result))
	}

	// Validate pagination parameters
	paginationResult := s.validator.ValidatePagination(page, limit)
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

	return s.repo.SearchNews(ctx, query, page, limit)
}
