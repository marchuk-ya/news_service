package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// News represents a news article in the system
type News struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title" validate:"required,min=3,max=200"`
	Content   string             `bson:"content" json:"content" validate:"required,min=10"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// NewsRepository defines the interface for news storage operations
type NewsRepository interface {
	Create(news *News) error
	GetByID(id string) (*News, error)
	GetAll(page, limit int) ([]*News, int64, error)
	Update(news *News) error
	Delete(id string) error
	Search(query string, page, limit int) ([]*News, int64, error)
}

// NewsService defines the interface for news business logic
type NewsService interface {
	CreateNews(news *News) error
	GetNewsByID(id string) (*News, error)
	GetAllNews(page, limit int) ([]*News, int64, error)
	UpdateNews(news *News) error
	DeleteNews(id string) error
	SearchNews(query string, page, limit int) ([]*News, int64, error)
}
