package repository

import (
	"context"

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

func (r mediaRepository) GetMedia(ctx context.Context, id string) (domain.Media, error) {
	filter := bson.D{{"_id", id}}
	media := domain.Media{}

	if err := r.db.FindOne(ctx, filter).Decode(&media); err != nil {
		if err == mongo.ErrNoDocuments {
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

func (r *mediaRepository) DeleteMedia(ctx context.Context, id string) error {
	find := bson.D{{"_id", id}}
	_, err := r.db.DeleteOne(ctx, find)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.ErrMediaNotFound
		}
		return err
	}

	return nil
}
