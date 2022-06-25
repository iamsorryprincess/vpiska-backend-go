package mongo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Events struct {
	db *mongo.Collection
}

func newEvents(db *mongo.Database, collectionName string) *Events {
	return &Events{
		db: db.Collection(collectionName),
	}
}

func (r *Events) CreateEvent(ctx context.Context, event domain.Event) (string, error) {
	event.ID = uuid.New().String()
	_, err := r.db.InsertOne(ctx, event)

	if err != nil {
		return "", err
	}

	return event.ID, nil
}

func (r *Events) GetEventById(ctx context.Context, id string) (domain.Event, error) {
	filter := bson.D{{Key: "_id", Value: id}, {Key: "state", Value: domain.EventStateOpened}}
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

func (r *Events) GetEventByOwnerId(ctx context.Context, ownerId string) (domain.Event, error) {
	filter := bson.D{{Key: "owner_id", Value: ownerId}, {Key: "state", Value: domain.EventStateOpened}}
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

func (r *Events) GetEventsByRange(ctx context.Context, xLeft float64, xRight float64, yLeft float64, yRight float64) ([]domain.EventRangeData, error) {
	filter := bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "state", Value: domain.EventStateOpened}},
		bson.D{{Key: "coordinates.x", Value: bson.D{{Key: "$gte", Value: xLeft}}}},
		bson.D{{Key: "coordinates.y", Value: bson.D{{Key: "$gte", Value: yLeft}}}},
		bson.D{{Key: "coordinates.x", Value: bson.D{{Key: "$lte", Value: xRight}}}},
		bson.D{{Key: "coordinates.y", Value: bson.D{{Key: "$lte", Value: yRight}}}},
	}}}

	cursor, err := r.db.Find(ctx, filter, options.Find().SetProjection(bson.D{
		{Key: "_id", Value: 1},
		{Key: "name", Value: 1},
		{Key: "coordinates", Value: 1},
		{Key: "users_count", Value: bson.D{{Key: "$size", Value: "$users"}}},
	}))

	if err != nil {
		return nil, err
	}

	result := make([]domain.EventRangeData, 0)
	err = cursor.All(ctx, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Events) UpdateEvent(ctx context.Context, id string, address string, coordinates domain.Coordinates) error {
	filter := bson.D{{Key: "_id", Value: id}, {Key: "state", Value: domain.EventStateOpened}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "address", Value: address},
		{Key: "coordinates", Value: coordinates},
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

func (r *Events) RemoveEvent(ctx context.Context, id string) error {
	filter := bson.D{{Key: "_id", Value: id}, {Key: "state", Value: domain.EventStateOpened}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "state", Value: domain.EventStateClosed}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *Events) AddMedia(ctx context.Context, id string, mediaInfo domain.MediaInfo) error {
	filter := bson.D{{Key: "_id", Value: id}, {Key: "state", Value: domain.EventStateOpened}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "media", Value: mediaInfo}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *Events) RemoveMedia(ctx context.Context, eventId string, mediaId string) error {
	filter := bson.D{{Key: "_id", Value: eventId}, {Key: "state", Value: domain.EventStateOpened}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "media", Value: bson.D{{Key: "_id", Value: mediaId}}}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *Events) AddUserInfo(ctx context.Context, eventId string, userInfo domain.UserInfo) error {
	filter := bson.D{{Key: "_id", Value: eventId}, {Key: "state", Value: domain.EventStateOpened}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "users", Value: userInfo}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

func (r *Events) RemoveUserInfo(ctx context.Context, eventId string, userId string) error {
	filter := bson.D{{Key: "_id", Value: eventId}, {Key: "state", Value: domain.EventStateOpened}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "users", Value: bson.D{{Key: "_id", Value: userId}}}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return err
}

func (r *Events) AddChatMessage(ctx context.Context, id string, chatMessage domain.ChatMessage) error {
	filter := bson.D{{Key: "_id", Value: id}, {Key: "state", Value: domain.EventStateOpened}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "chat_messages", Value: chatMessage}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}
