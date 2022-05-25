package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db *mongo.Collection
}

func newMongoUsers(db *mongo.Database, collectionName string) Users {
	return &userRepository{
		db: db.Collection(collectionName),
	}
}

func (r *userRepository) GetNamesCount(ctx context.Context, name string) (int64, error) {
	filter := bson.D{{"name", name}}
	count, err := r.db.CountDocuments(ctx, filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) GetPhonesCount(ctx context.Context, phone string) (int64, error) {
	filter := bson.D{{"phone", phone}}
	count, err := r.db.CountDocuments(ctx, filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user domain.User) (string, error) {
	user.ID = uuid.New().String()
	_, err := r.db.InsertOne(ctx, user)

	if err != nil {
		return "", err
	}

	return user.ID, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	filter := bson.D{{"_id", id}}
	model := domain.User{}

	if err := r.db.FindOne(ctx, filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model, domain.ErrUserNotFound
		}
		return model, err
	}

	return model, nil
}

func (r *userRepository) GetUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	filter := bson.D{{"phone", phone}}
	model := domain.User{}

	if err := r.db.FindOne(ctx, filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model, domain.ErrUserNotFound
		}
		return model, err
	}

	return model, nil
}

func (r *userRepository) ChangePassword(ctx context.Context, id string, password string) error {
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"password", password}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) SetImageId(ctx context.Context, userId string, imageId string) error {
	filter := bson.D{{"_id", userId}}
	update := bson.D{{"$set", bson.D{{"image_id", imageId}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UpdateName(ctx context.Context, userId string, name string) error {
	filter := bson.D{{"_id", userId}}
	update := bson.D{{"$set", bson.D{{"name", name}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UpdatePhone(ctx context.Context, userId string, phone string) error {
	filter := bson.D{{"_id", userId}}
	update := bson.D{{"$set", bson.D{{"phone", phone}}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UpdateNameAndPhone(ctx context.Context, userId string, name string, phone string) error {
	filter := bson.D{{"_id", userId}}
	update := bson.D{{"$set", bson.D{
		{"name", name},
		{"phone", phone},
	}}}
	result, err := r.db.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
