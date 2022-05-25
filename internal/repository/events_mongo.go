package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type eventsRepository struct {
	db *mongo.Collection
}

func newEventsMongo(db *mongo.Database, collectionName string) Events {
	return &eventsRepository{
		db: db.Collection(collectionName),
	}
}

func (r *eventsRepository) CreateEvent(ctx context.Context, event domain.Event) (string, error) {
	event.ID = uuid.New().String()
	_, err := r.db.InsertOne(ctx, event)

	if err != nil {
		return "", err
	}

	return event.ID, nil
}

func (r *eventsRepository) GetEventById(ctx context.Context, id string) (domain.Event, error) {
	filter := bson.D{{"_id", id}}
	event := domain.Event{}
	err := r.db.FindOne(ctx, filter).Decode(&event)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Event{}, domain.ErrEventNotFound
		}
		return domain.Event{}, err
	}

	return event, nil
}

func (r *eventsRepository) GetEventByOwnerId(ctx context.Context, ownerId string) (domain.Event, error) {
	filter := bson.D{{"owner_id", ownerId}}
	event := domain.Event{}
	err := r.db.FindOne(ctx, filter).Decode(&event)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Event{}, domain.ErrEventNotFound
		}
		return domain.Event{}, err
	}

	return event, nil
}

func (r *eventsRepository) UpdateEvent(ctx context.Context, id string, address string, coordinates domain.Coordinates) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{
		{"address", address},
		{"coordinates", coordinates},
	}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventsRepository) RemoveEvent(ctx context.Context, id string) error {
	filter := bson.D{{"_id", id}}
	result, err := r.db.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventsRepository) AddMedia(ctx context.Context, id string, mediaInfo domain.MediaInfo) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$push", bson.D{{"media", mediaInfo}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *eventsRepository) RemoveMedia(ctx context.Context, eventId string, mediaId string) error {
	filter := bson.D{{"_id", eventId}}
	update := bson.D{{"$pull", bson.D{{"media", bson.D{{"_id", mediaId}}}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}
