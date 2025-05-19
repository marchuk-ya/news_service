package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"news_service/internal/domain"
)

type MockNewsRepository struct {
	mock.Mock
}

func (m *MockNewsRepository) Create(news *domain.News) error {
	args := m.Called(news)
	return args.Error(0)
}

func (m *MockNewsRepository) GetByID(id string) (*domain.News, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.News), args.Error(1)
}

func (m *MockNewsRepository) GetAll(page, limit int) ([]*domain.News, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]*domain.News), args.Get(1).(int64), args.Error(2)
}

func (m *MockNewsRepository) Update(news *domain.News) error {
	args := m.Called(news)
	return args.Error(0)
}

func (m *MockNewsRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockNewsRepository) Search(query string, page, limit int) ([]*domain.News, int64, error) {
	args := m.Called(query, page, limit)
	return args.Get(0).([]*domain.News), args.Get(1).(int64), args.Error(2)
}

func TestNewsService_CreateNews(t *testing.T) {
	mockRepo := new(MockNewsRepository)
	service := NewNewsService(mockRepo)

	news := &domain.News{
		Title:   "Test News",
		Content: "Test Content",
	}

	mockRepo.On("Create", news).Return(nil)

	err := service.CreateNews(news)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNewsService_GetNewsByID(t *testing.T) {
	mockRepo := new(MockNewsRepository)
	service := NewNewsService(mockRepo)

	expectedNews := &domain.News{
		Title:   "Test News",
		Content: "Test Content",
	}

	mockRepo.On("GetByID", "test-id").Return(expectedNews, nil)

	news, err := service.GetNewsByID("test-id")
	assert.NoError(t, err)
	assert.Equal(t, expectedNews, news)
	mockRepo.AssertExpectations(t)
}

func TestNewsService_GetAllNews(t *testing.T) {
	mockRepo := new(MockNewsRepository)
	service := NewNewsService(mockRepo)

	expectedNews := []*domain.News{
		{Title: "News 1", Content: "Content 1"},
		{Title: "News 2", Content: "Content 2"},
	}

	mockRepo.On("GetAll", 1, 10).Return(expectedNews, int64(2), nil)

	news, total, err := service.GetAllNews(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, expectedNews, news)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestNewsService_UpdateNews(t *testing.T) {
	mockRepo := new(MockNewsRepository)
	service := NewNewsService(mockRepo)

	news := &domain.News{
		Title:   "Updated News",
		Content: "Updated Content",
	}

	mockRepo.On("Update", news).Return(nil)

	err := service.UpdateNews(news)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNewsService_DeleteNews(t *testing.T) {
	mockRepo := new(MockNewsRepository)
	service := NewNewsService(mockRepo)

	mockRepo.On("Delete", "test-id").Return(nil)

	err := service.DeleteNews("test-id")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNewsService_SearchNews(t *testing.T) {
	mockRepo := new(MockNewsRepository)
	service := NewNewsService(mockRepo)

	expectedNews := []*domain.News{
		{Title: "Golang News", Content: "Go programming language"},
	}

	mockRepo.On("Search", "golang", 1, 10).Return(expectedNews, int64(1), nil)

	news, total, err := service.SearchNews("golang", 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, expectedNews, news)
	assert.Equal(t, int64(1), total)
	mockRepo.AssertExpectations(t)
}
