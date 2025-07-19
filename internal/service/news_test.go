package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"news_service/internal/domain"
	"news_service/internal/mocks"
	"news_service/internal/tools"
)

func setupTestService(t *testing.T) (*mocks.NewsRepository, domain.NewsService) {
	mockRepo := mocks.NewNewsRepository(t)
	service := NewNewsService(mockRepo)
	return mockRepo, service
}

func TestGetAllNews(t *testing.T) {
	mockRepo, service := setupTestService(t)

	tests := []struct {
		name      string
		page      int
		limit     int
		mockNews  []*domain.News
		mockTotal int64
		mockErr   error
		wantErr   bool
	}{
		{
			name:  "successful pagination",
			page:  1,
			limit: 10,
			mockNews: []*domain.News{
				{
					ID:        primitive.NewObjectID(),
					Title:     "Test News 1",
					Content:   "Content 1",
					CreatedAt: tools.GetCurrentTime(),
					UpdatedAt: tools.GetCurrentTime(),
				},
				{
					ID:        primitive.NewObjectID(),
					Title:     "Test News 2",
					Content:   "Content 2",
					CreatedAt: tools.GetCurrentTime(),
					UpdatedAt: tools.GetCurrentTime(),
				},
			},
			mockTotal: 2,
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "empty result",
			page:      1,
			limit:     10,
			mockNews:  []*domain.News{},
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "repository error",
			page:      1,
			limit:     10,
			mockNews:  nil,
			mockTotal: 0,
			mockErr:   assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations using generated mock
			mockRepo.On("GetAll", mock.Anything, tt.page, tt.limit).
				Return(tt.mockNews, tt.mockTotal, tt.mockErr).Once()

			// Execute test
			news, total, err := service.GetAll(context.Background(), tt.page, tt.limit)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.mockErr, err)
				assert.Nil(t, news)
				assert.Equal(t, int64(0), total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockNews, news)
				assert.Equal(t, tt.mockTotal, total)
			}
		})
	}
}

func TestSearchNews(t *testing.T) {
	mockRepo, service := setupTestService(t)

	tests := []struct {
		name      string
		query     string
		page      int
		limit     int
		mockNews  []*domain.News
		mockTotal int64
		mockErr   error
		wantErr   bool
	}{
		{
			name:  "successful search",
			query: "test",
			page:  1,
			limit: 10,
			mockNews: []*domain.News{
				{
					ID:        primitive.NewObjectID(),
					Title:     "Test News",
					Content:   "Test Content",
					CreatedAt: tools.GetCurrentTime(),
					UpdatedAt: tools.GetCurrentTime(),
				},
			},
			mockTotal: 1,
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "no results",
			query:     "nonexistent",
			page:      1,
			limit:     10,
			mockNews:  []*domain.News{},
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "repository error",
			query:     "test",
			page:      1,
			limit:     10,
			mockNews:  nil,
			mockTotal: 0,
			mockErr:   assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations using generated mock
			mockRepo.On("SearchNews", mock.Anything, tt.query, tt.page, tt.limit).
				Return(tt.mockNews, tt.mockTotal, tt.mockErr).Once()

			news, total, err := service.SearchNews(context.Background(), tt.query, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.mockErr, err)
				assert.Nil(t, news)
				assert.Equal(t, int64(0), total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockNews, news)
				assert.Equal(t, tt.mockTotal, total)
			}
		})
	}
}

func TestCreateNews(t *testing.T) {
	mockRepo, service := setupTestService(t)

	tests := []struct {
		name    string
		news    *domain.News
		mockErr error
		wantErr bool
	}{
		{
			name: "successful creation",
			news: &domain.News{
				Title:   "Test News",
				Content: "Test Content",
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "repository error",
			news: &domain.News{
				Title:   "Test News",
				Content: "Test Content",
			},
			mockErr: assert.AnError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations using generated mock
			mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
				return n.Title == tt.news.Title && n.Content == tt.news.Content
			})).Return(tt.mockErr).Once()

			err := service.Create(context.Background(), tt.news)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.mockErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tt.news.CreatedAt)
				assert.NotZero(t, tt.news.UpdatedAt)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	mockRepo, service := setupTestService(t)

	tests := []struct {
		name     string
		id       string
		mockNews *domain.News
		mockErr  error
		wantErr  bool
	}{
		{
			name: "successful retrieval",
			id:   primitive.NewObjectID().Hex(),
			mockNews: &domain.News{
				ID:        primitive.NewObjectID(),
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:     "not found",
			id:       primitive.NewObjectID().Hex(),
			mockNews: nil,
			mockErr:  domain.ErrNotFound,
			wantErr:  true,
		},
		{
			name:     "empty id",
			id:       "",
			mockNews: nil,
			mockErr:  nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id != "" {
				// Setup mock expectations using generated mock
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, tt.mockErr).Once()
			}

			news, err := service.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.mockErr != nil {
					assert.ErrorIs(t, err, tt.mockErr)
				}
				assert.Nil(t, news)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockNews, news)
			}
		})
	}
}

func TestUpdateNews(t *testing.T) {
	mockRepo, service := setupTestService(t)

	tests := []struct {
		name     string
		news     *domain.News
		mockNews *domain.News
		mockErr  error
		wantErr  bool
	}{
		{
			name: "successful update",
			news: &domain.News{
				ID:      primitive.NewObjectID(),
				Title:   "Updated News",
				Content: "Updated Content",
			},
			mockNews: &domain.News{
				ID:        primitive.NewObjectID(),
				Title:     "Original News",
				Content:   "Original Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "news not found",
			news: &domain.News{
				ID:      primitive.NewObjectID(),
				Title:   "Updated News",
				Content: "Updated Content",
			},
			mockNews: nil,
			mockErr:  domain.ErrNotFound,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations using generated mock
			mockRepo.On("GetByID", mock.Anything, tt.news.ID.Hex()).
				Return(tt.mockNews, tt.mockErr).Once()

			if tt.mockErr == nil {
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
					return n.ID == tt.news.ID && n.Title == tt.news.Title && n.Content == tt.news.Content
				})).Return(nil).Once()
			}

			err := service.Update(context.Background(), tt.news)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.mockErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockNews.CreatedAt, tt.news.CreatedAt)
				assert.NotZero(t, tt.news.UpdatedAt)
			}
		})
	}
}

func TestDeleteNews(t *testing.T) {
	mockRepo, service := setupTestService(t)

	tests := []struct {
		name     string
		id       string
		mockNews *domain.News
		mockErr  error
		wantErr  bool
	}{
		{
			name: "successful deletion",
			id:   primitive.NewObjectID().Hex(),
			mockNews: &domain.News{
				ID:        primitive.NewObjectID(),
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:     "news not found",
			id:       primitive.NewObjectID().Hex(),
			mockNews: nil,
			mockErr:  domain.ErrNotFound,
			wantErr:  true,
		},
		{
			name:     "empty id",
			id:       "",
			mockNews: nil,
			mockErr:  nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id != "" {
				// Setup mock expectations using generated mock
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, tt.mockErr).Once()

				if tt.mockErr == nil {
					mockRepo.On("Delete", mock.Anything, tt.id).
						Return(nil).Once()
				}
			}

			err := service.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.mockErr != nil {
					assert.ErrorIs(t, err, tt.mockErr)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
