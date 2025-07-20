package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"news_service/internal/domain"
	"news_service/internal/mocks"
)

// Integration tests that test the interaction between layers
func TestNewsUseCaseIntegration(t *testing.T) {
	mockRepo := mocks.NewNewsRepository(t)
	useCase := NewNewsUseCase(mockRepo)

	t.Run("full CRUD workflow", func(t *testing.T) {
		ctx := context.Background()

		// Test Create
		title := "Integration Test News"
		content := "This is a test news article for integration testing"

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
			return n.Title == title && n.Content == content
		})).Run(func(args mock.Arguments) {
			// Set the ID on the news object to simulate repository behavior
			news := args.Get(1).(*domain.News)
			news.ID = primitive.NewObjectID().Hex()
		}).Return(nil).Once()

		news, err := useCase.CreateNews(ctx, title, content)
		require.NoError(t, err)
		require.NotNil(t, news)
		assert.Equal(t, title, news.Title)
		assert.Equal(t, content, news.Content)
		assert.NotEmpty(t, news.ID)
		assert.NotZero(t, news.CreatedAt)
		assert.NotZero(t, news.UpdatedAt)

		// Test GetByID
		mockRepo.On("GetByID", mock.Anything, news.ID).Return(news, nil).Once()

		retrievedNews, err := useCase.GetNewsByID(ctx, news.ID)
		require.NoError(t, err)
		assert.Equal(t, news, retrievedNews)

		// Test Update
		updatedTitle := "Updated Integration Test News"
		updatedContent := "This is an updated test news article"

		mockRepo.On("GetByID", mock.Anything, news.ID).Return(news, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(n *domain.News) bool {
			return n.Title == updatedTitle && n.Content == updatedContent
		})).Return(nil).Once()

		updatedNews, err := useCase.UpdateNews(ctx, news.ID, updatedTitle, updatedContent)
		require.NoError(t, err)
		assert.Equal(t, updatedTitle, updatedNews.Title)
		assert.Equal(t, updatedContent, updatedNews.Content)
		assert.NotZero(t, updatedNews.UpdatedAt)

		// Test Delete
		mockRepo.On("GetByID", mock.Anything, news.ID).Return(updatedNews, nil).Once()
		mockRepo.On("Delete", mock.Anything, news.ID).Return(nil).Once()

		err = useCase.DeleteNews(ctx, news.ID)
		require.NoError(t, err)
	})

	t.Run("pagination workflow", func(t *testing.T) {
		ctx := context.Background()

		mockNews := []*domain.News{
			{
				ID:        primitive.NewObjectID().Hex(),
				Title:     "News 1",
				Content:   "Content 1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        primitive.NewObjectID().Hex(),
				Title:     "News 2",
				Content:   "Content 2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockRepo.On("GetAll", mock.Anything, 1, 10).Return(mockNews, int64(2), nil).Once()

		news, total, err := useCase.GetAllNews(ctx, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, mockNews, news)
		assert.Equal(t, int64(2), total)
	})

	t.Run("search workflow", func(t *testing.T) {
		ctx := context.Background()
		query := "test"

		mockNews := []*domain.News{
			{
				ID:        primitive.NewObjectID().Hex(),
				Title:     "Test News",
				Content:   "Test Content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockRepo.On("SearchNews", mock.Anything, query, 1, 10).Return(mockNews, int64(1), nil).Once()

		news, total, err := useCase.SearchNews(ctx, query, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, mockNews, news)
		assert.Equal(t, int64(1), total)
	})
}

// Legacy service integration test removed - NewsService was deprecated in favor of NewsUseCase

func TestErrorHandlingIntegration(t *testing.T) {
	mockRepo := mocks.NewNewsRepository(t)
	useCase := NewNewsUseCase(mockRepo)

	t.Run("database errors are properly wrapped", func(t *testing.T) {
		ctx := context.Background()

		// Test Create with database error
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(domain.ErrDatabaseOperation).Once()

		news, err := useCase.CreateNews(ctx, "Test News", "Test Content")
		assert.Error(t, err)
		assert.Nil(t, news)
		assert.Contains(t, err.Error(), "database operation")

		// Test GetByID with not found error
		id := primitive.NewObjectID().Hex()
		mockRepo.On("GetByID", mock.Anything, id).Return(nil, domain.ErrNotFound).Once()

		news, err = useCase.GetNewsByID(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, news)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("validation errors are properly handled", func(t *testing.T) {
		ctx := context.Background()

		// Test Create with invalid input
		news, err := useCase.CreateNews(ctx, "", "Test Content")
		assert.Error(t, err)
		assert.Nil(t, news)
		assert.Contains(t, err.Error(), "validation")

		// Test Update with invalid input
		id := primitive.NewObjectID().Hex()
		news, err = useCase.UpdateNews(ctx, id, "", "Test Content")
		assert.Error(t, err)
		assert.Nil(t, news)
		assert.Contains(t, err.Error(), "validation")
	})
}

func TestPaginationIntegration(t *testing.T) {
	mockRepo := mocks.NewNewsRepository(t)
	useCase := NewNewsUseCase(mockRepo)

	t.Run("pagination defaults are applied", func(t *testing.T) {
		ctx := context.Background()

		// Test with invalid page and limit
		mockRepo.On("GetAll", mock.Anything, 1, 10).Return([]*domain.News{}, int64(0), nil).Once()

		news, total, err := useCase.GetAllNews(ctx, 0, 0)
		require.NoError(t, err)
		assert.Equal(t, []*domain.News{}, news)
		assert.Equal(t, int64(0), total)

		// Test with limit too high
		mockRepo.On("GetAll", mock.Anything, 1, 100).Return([]*domain.News{}, int64(0), nil).Once()

		news, total, err = useCase.GetAllNews(ctx, 1, 200)
		require.NoError(t, err)
		assert.Equal(t, []*domain.News{}, news)
		assert.Equal(t, int64(0), total)
	})
}

func TestSearchIntegration(t *testing.T) {
	mockRepo := mocks.NewNewsRepository(t)
	useCase := NewNewsUseCase(mockRepo)

	t.Run("search validation and pagination", func(t *testing.T) {
		ctx := context.Background()

		// Test with empty query
		news, total, err := useCase.SearchNews(ctx, "", 1, 10)
		assert.Error(t, err)
		assert.Nil(t, news)
		assert.Zero(t, total)

		// Test with valid query and invalid pagination
		query := "test"
		mockRepo.On("SearchNews", mock.Anything, query, 1, 10).Return([]*domain.News{}, int64(0), nil).Once()

		news, total, err = useCase.SearchNews(ctx, query, 0, 0)
		require.NoError(t, err)
		assert.Equal(t, []*domain.News{}, news)
		assert.Equal(t, int64(0), total)
	})
}
