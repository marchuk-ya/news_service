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
		{
			name:      "validation error - empty content",
			id:        primitive.NewObjectID().Hex(),
			title:     "Updated News",
			content:   "",
			mockNews:  nil,
			mockErr:   nil,
			wantErr:   true,
			wantTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock for each test
			mockRepo.ExpectedCalls = nil

			if tt.id != "" && !tt.wantErr {
				// Mock GetByID call for successful update
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, nil).Once()

				// Mock Update call
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
					return n.Title == tt.title && n.Content == tt.content
				})).Return(nil).Once()
			} else if tt.id != "" && tt.wantErr && tt.mockErr != nil {
				// Mock GetByID call for error cases
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
				assert.NotZero(t, news.UpdatedAt)
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
				// Mock GetByID call to check if news exists
				mockRepo.On("GetByID", mock.Anything, tt.id).
					Return(tt.mockNews, tt.mockErr).Once()

				if tt.mockErr == nil {
					// Mock Delete call for successful deletion
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
			// Reset mock for each test
			mockRepo.ExpectedCalls = nil

			if !tt.wantErr {
				mockRepo.On("SearchNews", mock.Anything, tt.query, tt.page, tt.limit).
					Return(tt.mockNews, tt.mockTotal, tt.mockErr).Once()
			} else if tt.mockErr != nil {
				// Mock for repository error case
				mockRepo.On("SearchNews", mock.Anything, tt.query, tt.page, tt.limit).
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

func TestGetAllNews(t *testing.T) {
	mockRepo, useCase := setupTestUseCase(t)

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
			name:  "successful retrieval",
			page:  1,
			limit: 10,
			mockNews: []*domain.News{
				{
					ID:        primitive.NewObjectID().Hex(),
					Title:     "Test News 1",
					Content:   "Test Content 1",
					CreatedAt: tools.GetCurrentTime(),
					UpdatedAt: tools.GetCurrentTime(),
				},
				{
					ID:        primitive.NewObjectID().Hex(),
					Title:     "Test News 2",
					Content:   "Test Content 2",
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
		{
			name:      "invalid pagination - negative page",
			page:      -1,
			limit:     10,
			mockNews:  []*domain.News{},
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   false, // Should apply defaults
		},
		{
			name:      "invalid pagination - zero limit",
			page:      1,
			limit:     0,
			mockNews:  []*domain.News{},
			mockTotal: 0,
			mockErr:   nil,
			wantErr:   false, // Should apply defaults
		},
		{
			name:      "invalid pagination - limit too high",
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
			// Reset mock for each test
			mockRepo.ExpectedCalls = nil

			// Determine expected page and limit after validation
			expectedPage := tt.page
			expectedLimit := tt.limit

			if expectedPage < 1 {
				expectedPage = 1
			}
			if expectedLimit < 1 {
				expectedLimit = 10
			}
			if expectedLimit > 100 {
				expectedLimit = 100
			}

			mockRepo.On("GetAll", mock.Anything, expectedPage, expectedLimit).
				Return(tt.mockNews, tt.mockTotal, tt.mockErr).Once()

			news, total, err := useCase.GetAllNews(context.Background(), tt.page, tt.limit)

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

// Benchmark tests
func BenchmarkCreateNewsUseCase(b *testing.B) {
	mockRepo := mocks.NewNewsRepository(nil)
	useCase := NewNewsUseCase(mockRepo)

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = useCase.CreateNews(context.Background(), "Benchmark News", "Benchmark Content")
	}
}

func BenchmarkGetNewsByIDUseCase(b *testing.B) {
	mockRepo := mocks.NewNewsRepository(nil)
	useCase := NewNewsUseCase(mockRepo)

	mockNews := &domain.News{
		ID:        primitive.NewObjectID().Hex(),
		Title:     "Benchmark News",
		Content:   "Benchmark Content",
		CreatedAt: tools.GetCurrentTime(),
		UpdatedAt: tools.GetCurrentTime(),
	}

	mockRepo.On("GetByID", mock.Anything, mock.Anything).Return(mockNews, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = useCase.GetNewsByID(context.Background(), mockNews.ID)
	}
}
