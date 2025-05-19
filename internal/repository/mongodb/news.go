package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"news_service/internal/domain"
)

const (
	collectionName = "news"
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

func (r *newsRepository) Create(news *domain.News) error {
	news.CreatedAt = time.Now()
	news.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(context.Background(), news)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		news.ID = oid
	}
	return nil
}

func (r *newsRepository) GetByID(id string) (*domain.News, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var news domain.News
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&news)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("news not found")
		}
		return nil, err
	}

	return &news, nil
}

func (r *newsRepository) GetAll(page, limit int) ([]*domain.News, int64, error) {
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var news []*domain.News
	if err = cursor.All(context.Background(), &news); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return news, total, nil
}

func (r *newsRepository) Update(news *domain.News) error {
	news.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":      news.Title,
			"content":    news.Content,
			"updated_at": news.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": news.ID},
		update,
	)
	return err
}

func (r *newsRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

func (r *newsRepository) Search(query string, page, limit int) ([]*domain.News, int64, error) {
	skip := (page - 1) * limit
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"content": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())

	var news []*domain.News
	if err = cursor.All(context.Background(), &news); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, 0, err
	}

	return news, total, nil
}
