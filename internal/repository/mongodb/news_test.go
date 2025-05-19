package mongodb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"news_service/internal/domain"
)

func setupTestDB(t *testing.T) (*mongo.Client, func()) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	// Clean up the test database
	database := client.Database("test_news_service")
	err = database.Drop(ctx)
	require.NoError(t, err)

	return client, func() {
		err := database.Drop(ctx)
		require.NoError(t, err)
		client.Disconnect(ctx)
	}
}

func TestNewsRepository_Create(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewNewsRepository(client, "test_news_service")

	news := &domain.News{
		Title:   "Test News",
		Content: "Test Content",
	}

	err := repo.Create(news)
	require.NoError(t, err)
	assert.NotEmpty(t, news.ID)
	assert.NotZero(t, news.CreatedAt)
	assert.NotZero(t, news.UpdatedAt)
}

func TestNewsRepository_GetByID(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewNewsRepository(client, "test_news_service")

	// Create test news
	news := &domain.News{
		Title:   "Test News",
		Content: "Test Content",
	}
	err := repo.Create(news)
	require.NoError(t, err)

	// Test getting the news
	retrieved, err := repo.GetByID(news.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, news.ID, retrieved.ID)
	assert.Equal(t, news.Title, retrieved.Title)
	assert.Equal(t, news.Content, retrieved.Content)
}

func TestNewsRepository_GetAll(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewNewsRepository(client, "test_news_service")

	// Create multiple test news
	for i := 0; i < 15; i++ {
		news := &domain.News{
			Title:   "Test News " + string(rune('A'+i)),
			Content: "Test Content " + string(rune('A'+i)),
		}
		err := repo.Create(news)
		require.NoError(t, err)
	}

	// Test getting all news with pagination
	news, total, err := repo.GetAll(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, news, 10)

	// Test second page
	news, total, err = repo.GetAll(2, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, news, 5)
}

func TestNewsRepository_Update(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewNewsRepository(client, "test_news_service")

	// Create test news
	news := &domain.News{
		Title:   "Test News",
		Content: "Test Content",
	}
	err := repo.Create(news)
	require.NoError(t, err)

	// Update the news
	news.Title = "Updated Title"
	news.Content = "Updated Content"
	err = repo.Update(news)
	require.NoError(t, err)

	// Verify the update
	retrieved, err := repo.GetByID(news.ID.Hex())
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", retrieved.Title)
	assert.Equal(t, "Updated Content", retrieved.Content)
}

func TestNewsRepository_Delete(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewNewsRepository(client, "test_news_service")

	// Create test news
	news := &domain.News{
		Title:   "Test News",
		Content: "Test Content",
	}
	err := repo.Create(news)
	require.NoError(t, err)

	// Delete the news
	err = repo.Delete(news.ID.Hex())
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(news.ID.Hex())
	assert.Error(t, err)
}

func TestNewsRepository_Search(t *testing.T) {
	client, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewNewsRepository(client, "test_news_service")

	// Create test news
	news := []*domain.News{
		{Title: "Golang News", Content: "Go programming language"},
		{Title: "Python News", Content: "Python programming language"},
		{Title: "Java News", Content: "Java programming language"},
	}

	for _, n := range news {
		err := repo.Create(n)
		require.NoError(t, err)
	}

	// Test search
	results, total, err := repo.Search("golang", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, results, 1)
	assert.Equal(t, "Golang News", results[0].Title)
}
