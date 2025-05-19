package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"news_service/internal/domain"
	"news_service/internal/handler"
	"news_service/internal/repository/mongodb"
	"news_service/internal/service"
)

func setupTestEnvironment(t *testing.T) (*gin.Engine, func()) {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	// Clean up the test database
	database := client.Database("test_news_service")
	err = database.Drop(ctx)
	require.NoError(t, err)

	// Initialize dependencies
	repo := mongodb.NewNewsRepository(client, "test_news_service")
	newsService := service.NewNewsService(repo)
	newsHandler := handler.NewNewsHandler(newsService)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	newsHandler.RegisterRoutes(router)

	// Return cleanup function
	cleanup := func() {
		err := database.Drop(ctx)
		require.NoError(t, err)
		client.Disconnect(ctx)
	}

	return router, cleanup
}

func TestNewsCRUD(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test Create
	news := &domain.News{
		Title:   "Test News",
		Content: "Test Content",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/news", nil)
	req.Form = map[string][]string{
		"title":   {news.Title},
		"content": {news.Content},
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusSeeOther, w.Code)

	// Test List
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test Search
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/news/search?q=Test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test Update
	updatedNews := &domain.News{
		Title:   "Updated News",
		Content: "Updated Content",
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/news/1", nil)
	req.Form = map[string][]string{
		"title":   {updatedNews.Title},
		"content": {updatedNews.Content},
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusSeeOther, w.Code)

	// Test Delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/news/1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewsValidation(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test Create with invalid data
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/news", nil)
	req.Form = map[string][]string{
		"title":   {"Te"}, // Too short
		"content": {"Co"}, // Too short
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test Update with invalid data
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/news/1", nil)
	req.Form = map[string][]string{
		"title":   {"Te"}, // Too short
		"content": {"Co"}, // Too short
	}
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewsPagination(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create multiple news items
	for i := 0; i < 15; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/news", nil)
		req.Form = map[string][]string{
			"title":   {string(rune('A' + i))},
			"content": {string(rune('A' + i))},
		}
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusSeeOther, w.Code)
	}

	// Test first page
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/?page=1&limit=10", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test second page
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/?page=2&limit=10", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewsSearch(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create test news
	news := []struct {
		title   string
		content string
	}{
		{"Golang News", "Go programming language"},
		{"Python News", "Python programming language"},
		{"Java News", "Java programming language"},
	}

	for _, n := range news {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/news", nil)
		req.Form = map[string][]string{
			"title":   {n.title},
			"content": {n.content},
		}
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusSeeOther, w.Code)
	}

	// Test search by title
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/news/search?q=golang", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test search by content
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/news/search?q=programming", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
