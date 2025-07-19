package domain

import (
	"context"
	"time"
)

// News represents a news article in the system
type News struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" validate:"required,min=3,max=200,notblank"`
	Content   string    `json:"content" validate:"required,min=10,notblank"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewsRepository defines the interface for news storage operations
type NewsRepository interface {
	Create(ctx context.Context, news *News) error
	GetByID(ctx context.Context, id string) (*News, error)
	GetAll(ctx context.Context, page, limit int) ([]*News, int64, error)
	Update(ctx context.Context, news *News) error
	Delete(ctx context.Context, id string) error
	SearchNews(ctx context.Context, query string, page, limit int) ([]*News, int64, error)
}

// NewsUseCase defines the interface for news business operations
type NewsUseCase interface {
	CreateNews(ctx context.Context, title, content string) (*News, error)
	GetNewsByID(ctx context.Context, id string) (*News, error)
	GetAllNews(ctx context.Context, page, limit int) ([]*News, int64, error)
	UpdateNews(ctx context.Context, id, title, content string) (*News, error)
	DeleteNews(ctx context.Context, id string) error
	SearchNews(ctx context.Context, query string, page, limit int) ([]*News, int64, error)
}

// NewsService defines the interface for news business operations (legacy - for backward compatibility)
type NewsService interface {
	Create(ctx context.Context, news *News) error
	GetByID(ctx context.Context, id string) (*News, error)
	GetAll(ctx context.Context, page, limit int) ([]*News, int64, error)
	Update(ctx context.Context, news *News) error
	Delete(ctx context.Context, id string) error
	SearchNews(ctx context.Context, query string, page, limit int) ([]*News, int64, error)
}
