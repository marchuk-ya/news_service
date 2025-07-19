package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// News represents a news article in the system
type News struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title" validate:"required,min=3,max=200,notblank" validateMsg:"Title is required, must be between 3 and 200 characters, and cannot be blank"`
	Content   string             `bson:"content" json:"content" validate:"required,min=10,notblank" validateMsg:"Content is required, must be at least 10 characters, and cannot be blank"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
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

// NewsService defines the interface for news business operations
type NewsService interface {
	Create(ctx context.Context, news *News) error
	GetByID(ctx context.Context, id string) (*News, error)
	GetAll(ctx context.Context, page, limit int) ([]*News, int64, error)
	Update(ctx context.Context, news *News) error
	Delete(ctx context.Context, id string) error
	SearchNews(ctx context.Context, query string, page, limit int) ([]*News, int64, error)
}
