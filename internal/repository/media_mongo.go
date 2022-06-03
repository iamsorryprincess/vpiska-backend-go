package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mediaRepository struct {
	db *mongo.Collection
}

func newMongoMedia(db *mongo.Database, collectionName string) Media {
	return &mediaRepository{
		db: db.Collection(collectionName),
	}
}

func (r *mediaRepository) GetMedia(ctx context.Context, id string) (domain.Media, error) {
	filter := bson.D{{"_id", id}}
	media := domain.Media{}

	if err := r.db.FindOne(ctx, filter).Decode(&media); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return media, domain.ErrMediaNotFound
		}
		return media, err
	}

	return media, nil
}

func (r *mediaRepository) CreateMedia(ctx context.Context, media domain.Media) (string, error) {
	media.ID = uuid.New().String()
	_, err := r.db.InsertOne(ctx, media)

	if err != nil {
		return "", err
	}

	return media.ID, nil
}

func (r *mediaRepository) UpdateMedia(ctx context.Context, media domain.Media) error {
	filter := bson.D{{"_id", media.ID}}
	update := bson.D{{"$set", bson.D{
		{"name", media.Name},
		{"content_type", media.ContentType},
		{"size", media.Size},
		{"last_modified_date", media.LastModifiedDate},
	}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrMediaNotFound
	}

	return nil
}

func (r *mediaRepository) DeleteMedia(ctx context.Context, id string) error {
	find := bson.D{{"_id", id}}
	result, err := r.db.DeleteOne(ctx, find)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrMediaNotFound
	}

	return nil
}
