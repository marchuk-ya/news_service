package repository_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"news_service/internal/domain"
	"news_service/internal/repository/mongodb"
	"news_service/tests/setup"
)

func TestNewsRepositoryEdgeCases(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setup.SetupTestMongoDB(t)
	defer cleanup()

	repo := mongodb.NewNewsRepository(client, "test_db")

	t.Run("Create with empty title", func(t *testing.T) {
		news := &domain.News{
			Title:     "",
			Content:   "Test content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, news)
		assert.Error(t, err, "Should fail with empty title")
	})

	t.Run("Create with empty content", func(t *testing.T) {
		news := &domain.News{
			Title:     "Test title",
			Content:   "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, news)
		assert.Error(t, err, "Should fail with empty content")
	})

	t.Run("GetByID with invalid ID format", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "invalid-id")
		assert.Error(t, err, "Should fail with invalid ID format")
	})

	t.Run("GetByID with non-existent ID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		_, err := repo.GetByID(ctx, nonExistentID)
		assert.Error(t, err, "Should fail with non-existent ID")
	})

	t.Run("GetAll with negative page", func(t *testing.T) {
		_, _, err := repo.GetAll(ctx, -1, 10)
		assert.Error(t, err, "Should fail with negative page")
	})

	t.Run("GetAll with zero limit", func(t *testing.T) {
		_, _, err := repo.GetAll(ctx, 1, 0)
		assert.Error(t, err, "Should fail with zero limit")
	})

	t.Run("GetAll with very large limit", func(t *testing.T) {
		_, _, err := repo.GetAll(ctx, 1, 1000000)
		assert.Error(t, err, "Should fail with very large limit")
	})

	t.Run("Update non-existent news", func(t *testing.T) {
		news := &domain.News{
			ID:        primitive.NewObjectID(),
			Title:     "Test title",
			Content:   "Test content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Update(ctx, news)
		assert.Error(t, err, "Should fail when updating non-existent news")
	})

	t.Run("Delete with invalid ID format", func(t *testing.T) {
		err := repo.Delete(ctx, "invalid-id")
		assert.Error(t, err, "Should fail with invalid ID format")
	})

	t.Run("Delete non-existent news", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		err := repo.Delete(ctx, nonExistentID)
		assert.Error(t, err, "Should fail when deleting non-existent news")
	})

	t.Run("SearchNews with empty query", func(t *testing.T) {
		_, _, err := repo.SearchNews(ctx, "", 1, 10)
		assert.Error(t, err, "Should fail with empty search query")
	})

	t.Run("SearchNews with special characters", func(t *testing.T) {
		_, _, err := repo.SearchNews(ctx, "!@#$%^&*()", 1, 10)
		assert.NoError(t, err, "Should handle special characters in search")
	})

	t.Run("SearchNews with very long query", func(t *testing.T) {
		longQuery := strings.Repeat("test query ", 100)
		_, _, err := repo.SearchNews(ctx, longQuery, 1, 10)
		assert.NoError(t, err, "Should handle very long search query")
	})

	t.Run("Concurrent operations", func(t *testing.T) {
		news := &domain.News{
			Title:     "Test title",
			Content:   "Test content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, news)
		require.NoError(t, err)

		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func() {
				_, err := repo.GetByID(ctx, news.ID.Hex())
				assert.NoError(t, err)

				news.Title = "Updated title"
				err = repo.Update(ctx, news)
				assert.NoError(t, err)

				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("Search with very long content", func(t *testing.T) {
		longContent := string(make([]byte, 1000000)) // 1MB of content
		news := &domain.News{
			Title:     "Test title",
			Content:   longContent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, news)
		require.NoError(t, err)

		_, _, err = repo.SearchNews(ctx, "test", 1, 10)
		assert.NoError(t, err, "Should handle search in very long content")
	})

	t.Run("Search with Unicode characters", func(t *testing.T) {
		news := &domain.News{
			Title:     "Тестовий заголовок",
			Content:   "Тестовий контент з Unicode символами",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, news)
		require.NoError(t, err)

		_, _, err = repo.SearchNews(ctx, "тест", 1, 10)
		assert.NoError(t, err, "Should handle search with Unicode characters")
	})
}
