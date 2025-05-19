package service

import (
	"news_service/internal/domain"
)

type newsService struct {
	repo domain.NewsRepository
}

// NewNewsService creates a new instance of news service
func NewNewsService(repo domain.NewsRepository) domain.NewsService {
	return &newsService{
		repo: repo,
	}
}

func (s *newsService) CreateNews(news *domain.News) error {
	return s.repo.Create(news)
}

func (s *newsService) GetNewsByID(id string) (*domain.News, error) {
	return s.repo.GetByID(id)
}

func (s *newsService) GetAllNews(page, limit int) ([]*domain.News, int64, error) {
	return s.repo.GetAll(page, limit)
}

func (s *newsService) UpdateNews(news *domain.News) error {
	return s.repo.Update(news)
}

func (s *newsService) DeleteNews(id string) error {
	return s.repo.Delete(id)
}

func (s *newsService) SearchNews(query string, page, limit int) ([]*domain.News, int64, error) {
	return s.repo.Search(query, page, limit)
}
