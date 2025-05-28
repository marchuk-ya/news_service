package mongodb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"news_service/internal/domain"
)

const (
	collectionName = "news"
	maxLimit       = 1000
)

var (
	ErrEmptyTitle   = errors.New("title cannot be empty")
	ErrEmptyContent = errors.New("content cannot be empty")
	ErrInvalidID    = errors.New("invalid ID format")
	ErrInvalidPage  = errors.New("page number must be positive")
	ErrInvalidLimit = errors.New("limit must be between 1 and 1000")
	ErrEmptyQuery   = errors.New("search query cannot be empty")
	ErrNewsNotFound = errors.New("news not found")
)

type newsRepository struct {
	client     *mongo.Client
	database   string
	collection *mongo.Collection
}

// NewNewsRepository creates a new instance of MongoDB news repository
func NewNewsRepository(client *mongo.Client, database string) domain.NewsRepository {
	collection := client.Database(database).Collection(collectionName)
	return &newsRepository{
		client:     client,
		database:   database,
		collection: collection,
	}
}

func (r *newsRepository) Create(ctx context.Context, news *domain.News) error {
	if news.Title == "" {
		return ErrEmptyTitle
	}
	if news.Content == "" {
		return ErrEmptyContent
	}

	news.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, news)
	return err
}

func (r *newsRepository) GetByID(ctx context.Context, id string) (*domain.News, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}

	var news domain.News
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&news)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNewsNotFound
		}
		return nil, err
	}

	return &news, nil
}

func (r *newsRepository) GetAll(ctx context.Context, page, limit int) ([]*domain.News, int64, error) {
	if page < 1 {
		return nil, 0, ErrInvalidPage
	}
	if limit < 1 || limit > maxLimit {
		return nil, 0, ErrInvalidLimit
	}

	skip := (page - 1) * limit

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var news []*domain.News
	if err = cursor.All(ctx, &news); err != nil {
		return nil, 0, err
	}

	return news, total, nil
}

func (r *newsRepository) Update(ctx context.Context, news *domain.News) error {
	if news.Title == "" {
		return ErrEmptyTitle
	}
	if news.Content == "" {
		return ErrEmptyContent
	}

	update := bson.M{
		"$set": bson.M{
			"title":      news.Title,
			"content":    news.Content,
			"updated_at": news.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": news.ID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNewsNotFound
	}
	return nil
}

func (r *newsRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNewsNotFound
	}
	return nil
}

func (r *newsRepository) SearchNews(ctx context.Context, query string, page, limit int) ([]*domain.News, int64, error) {
	if query == "" {
		return nil, 0, ErrEmptyQuery
	}
	if page < 1 {
		return nil, 0, ErrInvalidPage
	}
	if limit < 1 || limit > maxLimit {
		return nil, 0, ErrInvalidLimit
	}

	query = sanitizeQuery(query)
	if query == "" {
		return nil, 0, ErrEmptyQuery
	}

	query = escapeRegex(query)

	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"content": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find documents: %w", err)
	}
	defer cursor.Close(ctx)

	var news []*domain.News
	if err := cursor.All(ctx, &news); err != nil {
		return nil, 0, fmt.Errorf("failed to decode documents: %w", err)
	}

	return news, total, nil
}

// sanitizeQuery removes null bytes and other potentially dangerous characters from the query
func sanitizeQuery(query string) string {
	runes := []rune(query)

	var sanitized []rune
	for _, r := range runes {
		if r > 31 && r != 127 {
			sanitized = append(sanitized, r)
		}
	}

	if len(sanitized) > 1000 {
		sanitized = sanitized[:1000]
	}

	return string(sanitized)
}

// escapeRegex escapes special regex characters in the query
func escapeRegex(query string) string {
	specialChars := []string{`\`, `.`, `+`, `*`, `?`, `^`, `$`, `(`, `)`, `[`, `]`, `{`, `}`, `|`}
	result := query

	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, `\`+char)
	}

	return result
}

func (r *newsRepository) Search(ctx context.Context, query string) ([]*domain.News, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"content": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var news []*domain.News
	if err = cursor.All(ctx, &news); err != nil {
		return nil, err
	}

	return news, nil
}
