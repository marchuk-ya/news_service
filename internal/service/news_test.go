package service

import (
	"context"
	"testing"
	"time"

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

func setupTestUseCase(t *testing.T) (*mocks.NewsRepository, domain.NewsUseCase) {
	mockRepo := mocks.NewNewsRepository(t)
	useCase := NewNewsUseCase(mockRepo)
	return mockRepo, useCase
}

// Use Case Tests
func TestCreateNewsUseCase(t *testing.T) {
	mockRepo, useCase := setupTestUseCase(t)

	tests := []struct {
		name      string
		title     string
		content   string
		mockErr   error
		wantErr   bool
		wantTitle string
	}{
		{
			name:      "successful creation",
			title:     "Test News",
			content:   "Test Content",
			mockErr:   nil,
			wantErr:   false,
			wantTitle: "Test News",
		},
		{
			name:      "empty title",
			title:     "",
			content:   "Test Content",
			mockErr:   nil,
			wantErr:   true,
			wantTitle: "",
		},
		{
			name:      "empty content",
			title:     "Test News",
			content:   "",
			mockErr:   nil,
			wantErr:   true,
			wantTitle: "",
		},
		{
			name:      "repository error",
			title:     "Test News",
			content:   "Test Content",
			mockErr:   assert.AnError,
			wantErr:   true,
			wantTitle: "",
		},
		{
			name:      "title too short",
			title:     "AB",
			content:   "Test Content",
			mockErr:   nil,
			wantErr:   true,
			wantTitle: "",
		},
		{
			name:      "title too long",
			title:     string(make([]byte, 201)),
			content:   "Test Content",
			mockErr:   nil,
			wantErr:   true,
			wantTitle: "",
		},
		{
			name:      "content too short",
			title:     "Test News",
			content:   "Short",
			mockErr:   nil,
			wantErr:   true,
			wantTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock for each test
			mockRepo.ExpectedCalls = nil

			if !tt.wantErr {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
					return n.Title == tt.title && n.Content == tt.content
				})).Run(func(args mock.Arguments) {
					// Set the ID on the news object to simulate repository behavior
					news := args.Get(1).(*domain.News)
					news.ID = primitive.NewObjectID().Hex()
				}).Return(tt.mockErr).Once()
			} else if tt.mockErr != nil {
				// Mock for repository error case
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(tt.mockErr).Once()
			}

			news, err := useCase.CreateNews(context.Background(), tt.title, tt.content)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, news)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, news)
				assert.Equal(t, tt.wantTitle, news.Title)
				assert.NotZero(t, news.CreatedAt)
				assert.NotZero(t, news.UpdatedAt)
				assert.NotEmpty(t, news.ID)
			}
		})
	}
}

func TestGetNewsByIDUseCase(t *testing.T) {
	mockRepo, useCase := setupTestUseCase(t)

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
				ID:        primitive.NewObjectID().Hex(),
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:     "empty ID",
			id:       "",
			mockNews: nil,
			mockErr:  nil,
			wantErr:  true,
		},
		{
			name:     "not found",
			id:       primitive.NewObjectID().Hex(),
			mockNews: nil,
			mockErr:  domain.ErrNotFound,
			wantErr:  true,
		},
		{
			name:     "invalid ID format",
			id:       "invalid-id",
			mockNews: nil,
			mockErr:  domain.ErrInvalidInput,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id != "" {
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, tt.mockErr).Once()
			}

			news, err := useCase.GetNewsByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, news)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockNews, news)
			}
		})
	}
}

func TestUpdateNewsUseCase(t *testing.T) {
	mockRepo, useCase := setupTestUseCase(t)

	tests := []struct {
		name      string
		id        string
		title     string
		content   string
		mockNews  *domain.News
		mockErr   error
		wantErr   bool
		wantTitle string
	}{
		{
			name:    "successful update",
			id:      primitive.NewObjectID().Hex(),
			title:   "Updated News",
			content: "Updated Content",
			mockNews: &domain.News{
				ID:        primitive.NewObjectID().Hex(),
				Title:     "Original News",
				Content:   "Original Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr:   nil,
			wantErr:   false,
			wantTitle: "Updated News",
		},
		{
			name:      "empty ID",
			id:        "",
			title:     "Updated News",
			content:   "Updated Content",
			mockNews:  nil,
			mockErr:   nil,
			wantErr:   true,
			wantTitle: "",
		},
		{
			name:      "not found",
			id:        primitive.NewObjectID().Hex(),
			title:     "Updated News",
			content:   "Updated Content",
			mockNews:  nil,
			mockErr:   domain.ErrNotFound,
			wantErr:   true,
			wantTitle: "",
		},
		{
			name:      "validation error - empty title",
			id:        primitive.NewObjectID().Hex(),
			title:     "",
			content:   "Updated Content",
			mockNews:  nil, // не мокати GetByID
			mockErr:   nil,
			wantErr:   true,
			wantTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock for each test
			mockRepo.ExpectedCalls = nil

			// Додаємо мок GetByID лише якщо title не порожній і не очікується помилка валідації
			if tt.id != "" && !tt.wantErr && tt.title != "" {
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, tt.mockErr).Once()
				if tt.mockNews != nil {
					mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
						return n.Title == tt.title && n.Content == tt.content
					})).Return(nil).Once()
				}
			} else if tt.id != "" && tt.mockErr != nil {
				// Для not found
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, tt.mockErr).Once()
			}

			news, err := useCase.UpdateNews(context.Background(), tt.id, tt.title, tt.content)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, news)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, news)
				assert.Equal(t, tt.wantTitle, news.Title)
			}
		})
	}
}

func TestDeleteNewsUseCase(t *testing.T) {
	mockRepo, useCase := setupTestUseCase(t)

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
				ID:        primitive.NewObjectID().Hex(),
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:     "empty ID",
			id:       "",
			mockNews: nil,
			mockErr:  nil,
			wantErr:  true,
		},
		{
			name:     "not found",
			id:       primitive.NewObjectID().Hex(),
			mockNews: nil,
			mockErr:  domain.ErrNotFound,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id != "" {
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, tt.mockErr).Once()

				if tt.mockNews != nil {
					mockRepo.On("Delete", mock.Anything, tt.id).
						Return(nil).Once()
				}
			}

			err := useCase.DeleteNews(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearchNewsUseCase(t *testing.T) {
	mockRepo, useCase := setupTestUseCase(t)

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
					ID:        primitive.NewObjectID().Hex(),
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
			name:      "empty query",
			query:     "",
			page:      1,
			limit:     10,
			mockNews:  nil,
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   true,
		},
		{
			name:      "invalid pagination",
			query:     "test",
			page:      0,
			limit:     0,
			mockNews:  nil,
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   false, // Should apply defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockRepo.On("SearchNews", mock.Anything, tt.query, mock.Anything, mock.Anything).
					Return(tt.mockNews, tt.mockTotal, tt.mockErr).Once()
			}

			news, total, err := useCase.SearchNews(context.Background(), tt.query, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, news)
				assert.Zero(t, total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockNews, news)
				assert.Equal(t, tt.mockTotal, total)
			}
		})
	}
}

// Legacy Service Tests
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
					ID:        primitive.NewObjectID().Hex(),
					Title:     "Test News 1",
					Content:   "Content 1",
					CreatedAt: tools.GetCurrentTime(),
					UpdatedAt: tools.GetCurrentTime(),
				},
				{
					ID:        primitive.NewObjectID().Hex(),
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
			name:      "invalid page",
			page:      0,
			limit:     10,
			mockNews:  []*domain.News{},
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   false, // Should apply defaults
		},
		{
			name:      "invalid limit",
			page:      1,
			limit:     0,
			mockNews:  []*domain.News{},
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   false, // Should apply defaults
		},
		{
			name:      "limit too high",
			page:      1,
			limit:     200,
			mockNews:  []*domain.News{},
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   false, // Should cap at 100
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
				Return(tt.mockNews, tt.mockTotal, tt.mockErr).Once()

			news, total, err := service.GetAll(context.Background(), tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, news)
				assert.Zero(t, total)
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
					ID:        primitive.NewObjectID().Hex(),
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
			name:      "empty query",
			query:     "",
			page:      1,
			limit:     10,
			mockNews:  nil,
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockRepo.On("SearchNews", mock.Anything, tt.query, mock.Anything, mock.Anything).
					Return(tt.mockNews, tt.mockTotal, tt.mockErr).Once()
			}

			news, total, err := service.SearchNews(context.Background(), tt.query, tt.page, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, news)
				assert.Zero(t, total)
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
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "validation error",
			news: &domain.News{
				Title:     "",
				Content:   "Test Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
					return n.Title == tt.news.Title && n.Content == tt.news.Content
				})).Run(func(args mock.Arguments) {
					// Set the ID on the news object to simulate repository behavior
					news := args.Get(1).(*domain.News)
					news.ID = primitive.NewObjectID().Hex()
				}).Return(tt.mockErr).Once()
			}

			err := service.Create(context.Background(), tt.news)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.news.ID)
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
				ID:        primitive.NewObjectID().Hex(),
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name:     "empty ID",
			id:       "",
			mockNews: nil,
			mockErr:  nil,
			wantErr:  true,
		},
		{
			name:     "not found",
			id:       primitive.NewObjectID().Hex(),
			mockNews: nil,
			mockErr:  domain.ErrNotFound,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id != "" {
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, tt.mockErr).Once()
			}

			news, err := service.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
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
		name    string
		news    *domain.News
		mockErr error
		wantErr bool
	}{
		{
			name: "successful update",
			news: &domain.News{
				ID:        primitive.NewObjectID().Hex(),
				Title:     "Updated News",
				Content:   "Updated Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "validation error",
			news: &domain.News{
				ID:        primitive.NewObjectID().Hex(),
				Title:     "",
				Content:   "Updated Content",
				CreatedAt: tools.GetCurrentTime(),
				UpdatedAt: tools.GetCurrentTime(),
			},
			mockErr: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				// Mock GetByID call that happens in UpdateNews
				mockRepo.On("GetByID", mock.Anything, tt.news.ID).Return(tt.news, nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
					return n.Title == tt.news.Title && n.Content == tt.news.Content
				})).Return(tt.mockErr).Once()
			}

			err := service.Update(context.Background(), tt.news)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteNews(t *testing.T) {
	mockRepo, service := setupTestService(t)

	tests := []struct {
		name    string
		id      string
		mockErr error
		wantErr bool
	}{
		{
			name:    "successful deletion",
			id:      primitive.NewObjectID().Hex(),
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "empty ID",
			id:      "",
			mockErr: nil,
			wantErr: true,
		},
		{
			name:    "not found",
			id:      primitive.NewObjectID().Hex(),
			mockErr: domain.ErrNotFound,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id != "" {
				// Mock GetByID call that happens in DeleteNews
				mockRepo.On("GetByID", mock.Anything, tt.id).Return(&domain.News{ID: tt.id}, nil).Once()
				mockRepo.On("Delete", mock.Anything, tt.id).
					Return(tt.mockErr).Once()
			}

			err := service.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkCreateNewsUseCase(b *testing.B) {
	mockRepo := mocks.NewNewsRepository(nil)
	useCase := NewNewsUseCase(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = useCase.CreateNews(context.Background(), "Test News", "Test Content")
	}
}

func BenchmarkGetNewsByIDUseCase(b *testing.B) {
	mockRepo := mocks.NewNewsRepository(nil)
	useCase := NewNewsUseCase(mockRepo)

	mockNews := &domain.News{
		ID:        primitive.NewObjectID().Hex(),
		Title:     "Test News",
		Content:   "Test Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(mockNews, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = useCase.GetNewsByID(context.Background(), mockNews.ID)
	}
}
