package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"news_service/internal/domain"
	"news_service/internal/mocks"
	"news_service/internal/service"
)

func setupTestAPIHandler(t *testing.T) (*gin.Engine, *mocks.NewsRepository, *NewsAPIHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockRepo := mocks.NewNewsRepository(t)
	newsUseCase := service.NewNewsUseCase(mockRepo)
	handler := NewNewsAPIHandler(newsUseCase)

	handler.RegisterAPIRoutes(router)

	return router, mockRepo, handler
}

func TestAPIListNews(t *testing.T) {
	router, mockRepo, _ := setupTestAPIHandler(t)

	tests := []struct {
		name           string
		page           string
		limit          string
		mockNews       []*domain.News
		mockTotal      int64
		expectedStatus int
		expectedTotal  int64
	}{
		{
			name:  "successful pagination",
			page:  "1",
			limit: "10",
			mockNews: []*domain.News{
				{
					ID:        primitive.NewObjectID().Hex(),
					Title:     "Test News 1",
					Content:   "Content 1",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        primitive.NewObjectID().Hex(),
					Title:     "Test News 2",
					Content:   "Content 2",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockTotal:      2,
			expectedStatus: http.StatusOK,
			expectedTotal:  2,
		},
		{
			name:           "empty result",
			page:           "1",
			limit:          "10",
			mockNews:       []*domain.News{},
			mockTotal:      0,
			expectedStatus: http.StatusOK,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations using generated mock
			mockRepo.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
				Return(tt.mockNews, tt.mockTotal, nil).Once()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/news?page="+tt.page+"&limit="+tt.limit, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response
			var response ListNewsResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTotal, response.Total)
			assert.Len(t, response.News, len(tt.mockNews))
		})
	}
}

func TestAPICreateNews(t *testing.T) {
	router, mockRepo, _ := setupTestAPIHandler(t)

	tests := []struct {
		name           string
		request        CreateNewsRequest
		mockErr        error
		expectedStatus int
	}{
		{
			name: "successful creation",
			request: CreateNewsRequest{
				Title:   "Test News",
				Content: "Test Content",
			},
			mockErr:        nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "validation error - empty title",
			request: CreateNewsRequest{
				Title:   "",
				Content: "Test Content",
			},
			mockErr:        nil,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatus == http.StatusCreated {
				// Setup mock expectations using generated mock
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
					return n.Title == tt.request.Title && n.Content == tt.request.Content
				})).Run(func(args mock.Arguments) {
					// Set the ID on the news object to simulate repository behavior
					news := args.Get(1).(*domain.News)
					news.ID = primitive.NewObjectID().Hex()
				}).Return(tt.mockErr).Once()
			}

			requestBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/news", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAPIGetNews(t *testing.T) {
	router, mockRepo, _ := setupTestAPIHandler(t)

	tests := []struct {
		name           string
		id             string
		mockNews       *domain.News
		mockErr        error
		expectedStatus int
	}{
		{
			name: "successful retrieval",
			id:   primitive.NewObjectID().Hex(),
			mockNews: &domain.News{
				ID:        primitive.NewObjectID().Hex(),
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found",
			id:             primitive.NewObjectID().Hex(),
			mockNews:       nil,
			mockErr:        domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations using generated mock
			mockRepo.On("GetByID", mock.Anything, tt.id).
				Return(tt.mockNews, tt.mockErr).Once()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/news/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response NewsResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockNews.Title, response.Title)
				assert.Equal(t, tt.mockNews.Content, response.Content)
			}
		})
	}
}

func TestAPIUpdateNews(t *testing.T) {
	// Use fixed ID for consistent testing
	fixedID := primitive.NewObjectID()
	fixedIDHex := fixedID.Hex()

	tests := []struct {
		name           string
		id             string
		request        UpdateNewsRequest
		mockNews       *domain.News
		mockErr        error
		expectedStatus int
	}{
		{
			name: "successful update",
			id:   fixedIDHex,
			request: UpdateNewsRequest{
				Title:   "Updated News",
				Content: "Updated Content",
			},
			mockNews: &domain.News{
				ID:        fixedIDHex,
				Title:     "Original News",
				Content:   "Original Content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name: "not found",
			id:   primitive.NewObjectID().Hex(),
			request: UpdateNewsRequest{
				Title:   "Updated News",
				Content: "Updated Content",
			},
			mockNews:       nil,
			mockErr:        domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockRepo, _ := setupTestAPIHandler(t)

			if tt.expectedStatus == http.StatusOK {
				// Mock GetByID call for successful update
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, nil).Once()

				// Mock Update call
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
					return n.Title == tt.request.Title && n.Content == tt.request.Content
				})).Return(nil).Once()
			} else {
				// Mock GetByID call for error cases
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, tt.mockErr).Once()
			}

			requestBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/news/"+tt.id, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response NewsResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.request.Title, response.Title)
				assert.Equal(t, tt.request.Content, response.Content)
			}
		})
	}
}

func TestAPIDeleteNews(t *testing.T) {
	router, mockRepo, _ := setupTestAPIHandler(t)

	tests := []struct {
		name           string
		id             string
		mockNews       *domain.News
		mockErr        error
		expectedStatus int
	}{
		{
			name: "successful deletion",
			id:   primitive.NewObjectID().Hex(),
			mockNews: &domain.News{
				ID:        primitive.NewObjectID().Hex(),
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockErr:        nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "not found",
			id:             primitive.NewObjectID().Hex(),
			mockNews:       nil,
			mockErr:        domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock GetByID call to check if news exists
			mockRepo.On("GetByID", mock.Anything, tt.id).
				Return(tt.mockNews, tt.mockErr).Once()

			if tt.mockErr == nil {
				// Mock Delete call for successful deletion
				mockRepo.On("Delete", mock.Anything, tt.id).
					Return(nil).Once()
			}

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/news/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAPISearchNews(t *testing.T) {
	router, mockRepo, _ := setupTestAPIHandler(t)

	tests := []struct {
		name           string
		query          string
		page           string
		limit          string
		mockNews       []*domain.News
		mockTotal      int64
		expectedStatus int
		expectedTotal  int64
	}{
		{
			name:  "successful search",
			query: "test",
			page:  "1",
			limit: "10",
			mockNews: []*domain.News{
				{
					ID:        primitive.NewObjectID().Hex(),
					Title:     "Test News",
					Content:   "Test Content",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			mockTotal:      1,
			expectedStatus: http.StatusOK,
			expectedTotal:  1,
		},
		{
			name:           "empty query",
			query:          "",
			page:           "1",
			limit:          "10",
			mockNews:       nil,
			mockTotal:      0,
			expectedStatus: http.StatusBadRequest,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatus == http.StatusOK {
				// Setup mock expectations using generated mock
				mockRepo.On("SearchNews", mock.Anything, tt.query, mock.Anything, mock.Anything).
					Return(tt.mockNews, tt.mockTotal, nil).Once()
			}

			req := httptest.NewRequest(http.MethodGet, "/api/v1/news/search?q="+tt.query+"&page="+tt.page+"&limit="+tt.limit, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response SearchNewsResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTotal, response.Total)
				assert.Equal(t, tt.query, response.Query)
				assert.Len(t, response.News, len(tt.mockNews))
			}
		})
	}
}
