package service

import (
	"context"

	"news_service/internal/domain"
	"news_service/internal/tools"
)

type newsService struct {
	repo domain.NewsRepository
}

func NewNewsService(repo domain.NewsRepository) domain.NewsRepository {
	return &newsService{
		repo: repo,
	}
}

func (s *newsService) Create(ctx context.Context, news *domain.News) error {
	now := tools.GetCurrentTime()
	news.CreatedAt = now
	news.UpdatedAt = now
	return s.repo.Create(ctx, news)
}

func (s *newsService) GetByID(ctx context.Context, id string) (*domain.News, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *newsService) GetAll(ctx context.Context, page, limit int) ([]*domain.News, int64, error) {
	return s.repo.GetAll(ctx, page, limit)
}

func (s *newsService) Update(ctx context.Context, news *domain.News) error {
	news.UpdatedAt = tools.GetCurrentTime()
	return s.repo.Update(ctx, news)
}

func (s *newsService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *newsService) SearchNews(ctx context.Context, query string, page, limit int) ([]*domain.News, int64, error) {
	return s.repo.SearchNews(ctx, query, page, limit)
}
