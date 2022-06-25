package mongo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Media struct {
	db *mongo.Collection
}

func newMedia(db *mongo.Database, collectionName string) *Media {
	return &Media{
		db: db.Collection(collectionName),
	}
}

func (r *Media) GetMedia(ctx context.Context, id string) (domain.Media, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	media := domain.Media{}

	if err := r.db.FindOne(ctx, filter).Decode(&media); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return media, domain.ErrMediaNotFound
		}
		return media, err
	}

	return media, nil
}

func (r *Media) CreateMedia(ctx context.Context, media domain.Media) (string, error) {
	media.ID = uuid.New().String()
	_, err := r.db.InsertOne(ctx, media)

	if err != nil {
		return "", err
	}

	return media.ID, nil
}

func (r *Media) UpdateMedia(ctx context.Context, media domain.Media) error {
	filter := bson.D{{Key: "_id", Value: media.ID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: media.Name},
		{Key: "content_type", Value: media.ContentType},
		{Key: "size", Value: media.Size},
		{Key: "last_modified_date", Value: media.LastModifiedDate},
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

func (r *Media) DeleteMedia(ctx context.Context, id string) error {
	find := bson.D{{Key: "_id", Value: id}}
	result, err := r.db.DeleteOne(ctx, find)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrMediaNotFound
	}

	return nil
}
