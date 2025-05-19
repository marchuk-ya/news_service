package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"news_service/internal/domain"
)

type MockNewsService struct {
	mock.Mock
}

func (m *MockNewsService) CreateNews(news *domain.News) error {
	args := m.Called(news)
	return args.Error(0)
}

func (m *MockNewsService) GetNewsByID(id string) (*domain.News, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.News), args.Error(1)
}

func (m *MockNewsService) GetAllNews(page, limit int) ([]*domain.News, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]*domain.News), args.Get(1).(int64), args.Error(2)
}

func (m *MockNewsService) UpdateNews(news *domain.News) error {
	args := m.Called(news)
	return args.Error(0)
}

func (m *MockNewsService) DeleteNews(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockNewsService) SearchNews(query string, page, limit int) ([]*domain.News, int64, error) {
	args := m.Called(query, page, limit)
	return args.Get(0).([]*domain.News), args.Get(1).(int64), args.Error(2)
}

func setupTestRouter(service domain.NewsService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewNewsHandler(service)
	handler.RegisterRoutes(router)
	return router
}

func TestNewsHandler_ListNews(t *testing.T) {
	mockService := new(MockNewsService)
	router := setupTestRouter(mockService)

	expectedNews := []*domain.News{
		{Title: "News 1", Content: "Content 1"},
		{Title: "News 2", Content: "Content 2"},
	}

	mockService.On("GetAllNews", 1, 10).Return(expectedNews, int64(2), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestNewsHandler_CreateNews(t *testing.T) {
	mockService := new(MockNewsService)
	router := setupTestRouter(mockService)

	news := &domain.News{
		Title:   "Test News",
		Content: "Test Content",
	}

	mockService.On("CreateNews", mock.AnythingOfType("*domain.News")).Return(nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(news)
	req, _ := http.NewRequest("POST", "/news", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	mockService.AssertExpectations(t)
}

func TestNewsHandler_GetNews(t *testing.T) {
	mockService := new(MockNewsService)
	router := setupTestRouter(mockService)

	expectedNews := &domain.News{
		Title:   "Test News",
		Content: "Test Content",
	}

	mockService.On("GetNewsByID", "test-id").Return(expectedNews, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/news/test-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestNewsHandler_UpdateNews(t *testing.T) {
	mockService := new(MockNewsService)
	router := setupTestRouter(mockService)

	news := &domain.News{
		Title:   "Updated News",
		Content: "Updated Content",
	}

	mockService.On("GetNewsByID", "test-id").Return(news, nil)
	mockService.On("UpdateNews", mock.AnythingOfType("*domain.News")).Return(nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(news)
	req, _ := http.NewRequest("PUT", "/news/test-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	mockService.AssertExpectations(t)
}

func TestNewsHandler_DeleteNews(t *testing.T) {
	mockService := new(MockNewsService)
	router := setupTestRouter(mockService)

	mockService.On("DeleteNews", "test-id").Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/news/test-id", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestNewsHandler_SearchNews(t *testing.T) {
	mockService := new(MockNewsService)
	router := setupTestRouter(mockService)

	expectedNews := []*domain.News{
		{Title: "Golang News", Content: "Go programming language"},
	}

	mockService.On("SearchNews", "golang", 1, 10).Return(expectedNews, int64(1), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/news/search?q=golang", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
