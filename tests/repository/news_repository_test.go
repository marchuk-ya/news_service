package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"news_service/internal/domain"
	"news_service/internal/repository/mongodb"
	"news_service/tests/setup"
)

func TestNewsRepository(t *testing.T) {
	ctx := context.Background()
	client, cleanup := setup.SetupTestMongoDB(t)
	defer cleanup()

	repo := mongodb.NewNewsRepository(client, "test_db")

	news := &domain.News{
		Title:     "Test Title",
		Content:   "Test Content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("Create and Get", func(t *testing.T) {
		err := repo.Create(ctx, news)
		require.NoError(t, err)
		assert.NotEmpty(t, news.ID)

		retrieved, err := repo.GetByID(ctx, news.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, news.Title, retrieved.Title)
		assert.Equal(t, news.Content, retrieved.Content)
	})

	t.Run("GetAll with pagination", func(t *testing.T) {
		for i := 0; i < 15; i++ {
			item := &domain.News{
				Title:     "Test Title " + string(rune('A'+i)),
				Content:   "Test Content " + string(rune('A'+i)),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := repo.Create(ctx, item)
			require.NoError(t, err)
		}

		items, total, err := repo.GetAll(ctx, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(16), total)
		assert.Len(t, items, 10)

		items, total, err = repo.GetAll(ctx, 2, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(16), total)
		assert.Len(t, items, 6)
	})

	t.Run("Update", func(t *testing.T) {
		news.Title = "Updated Title"
		news.Content = "Updated Content"
		news.UpdatedAt = time.Now()

		err := repo.Update(ctx, news)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(ctx, news.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, "Updated Title", retrieved.Title)
		assert.Equal(t, "Updated Content", retrieved.Content)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, news.ID.Hex())
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, news.ID.Hex())
		assert.Error(t, err)
	})

	t.Run("Search", func(t *testing.T) {
		testNews := []*domain.News{
			{
				Title:     "Golang Programming",
				Content:   "Learn Go programming language",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Title:     "Python Programming",
				Content:   "Learn Python programming language",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		for _, item := range testNews {
			err := repo.Create(ctx, item)
			require.NoError(t, err)
		}

		results, total, err := repo.SearchNews(ctx, "programming", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, results, 2)
	})
}
