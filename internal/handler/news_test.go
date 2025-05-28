package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"news_service/internal/domain"
	"news_service/internal/mocks"
	"news_service/internal/tools"
)

type mockTemplate struct {
	name string
	data interface{}
}

func (m *mockTemplate) Render(w http.ResponseWriter) error {
	return nil
}

func (m *mockTemplate) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (m *mockTemplate) Instance(name string, data interface{}) render.Render {
	return &mockTemplate{
		name: name,
		data: data,
	}
}

func setupTestHandler() (*gin.Engine, *mocks.NewsRepository, *NewsHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.HTMLRender = &mockTemplate{}

	mockRepo := new(mocks.NewsRepository)
	handler := NewNewsHandler(mockRepo)

	handler.RegisterRoutes(router)

	return router, mockRepo, handler
}

func TestListNews(t *testing.T) {
	router, mockRepo, _ := setupTestHandler()

	tests := []struct {
		name           string
		page           string
		limit          string
		mockNews       []*domain.News
		mockTotal      int64
		expectedStatus int
	}{
		{
			name:  "successful pagination",
			page:  "1",
			limit: "10",
			mockNews: []*domain.News{
				{
					ID:        primitive.NewObjectID(),
					Title:     "Test News 1",
					Content:   "Content 1",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        primitive.NewObjectID(),
					Title:     "Test News 2",
					Content:   "Content 2",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockTotal:      2,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty result",
			page:           "1",
			limit:          "10",
			mockNews:       []*domain.News{},
			mockTotal:      0,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid page",
			page:           "0",
			limit:          "10",
			mockNews:       []*domain.News{},
			mockTotal:      0,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
				Return(tt.mockNews, tt.mockTotal, nil).Once()

			req := httptest.NewRequest(http.MethodGet, "/?page="+tt.page+"&limit="+tt.limit, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSearchNews(t *testing.T) {
	router, mockRepo, _ := setupTestHandler()

	tests := []struct {
		name           string
		query          string
		page           string
		limit          string
		mockNews       []*domain.News
		mockTotal      int64
		expectedStatus int
	}{
		{
			name:  "successful search",
			query: "test",
			page:  "1",
			limit: "10",
			mockNews: []*domain.News{
				{
					ID:        primitive.NewObjectID(),
					Title:     "Test News",
					Content:   "Test Content",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockTotal:      1,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "no results",
			query:          "nonexistent",
			page:           "1",
			limit:          "10",
			mockNews:       []*domain.News{},
			mockTotal:      0,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("SearchNews", mock.Anything, tt.query, mock.Anything, mock.Anything).
				Return(tt.mockNews, tt.mockTotal, nil).Once()

			req := httptest.NewRequest(http.MethodGet, "/news/search?q="+tt.query+"&page="+tt.page+"&limit="+tt.limit, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateNews(t *testing.T) {
	router, mockRepo, _ := setupTestHandler()

	tests := []struct {
		name           string
		news           domain.News
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful creation",
			news: domain.News{
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockError:      nil,
			expectedStatus: http.StatusSeeOther,
		},
		{
			name: "validation error",
			news: domain.News{
				Title:     "",
				Content:   "Test Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockError == nil && tt.expectedStatus == http.StatusSeeOther {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
					return n.Title == tt.news.Title && n.Content == tt.news.Content
				})).Return(nil).Once()
			}

			formData := strings.NewReader("title=" + tt.news.Title + "&content=" + tt.news.Content)
			req := httptest.NewRequest(http.MethodPost, "/news", formData)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
