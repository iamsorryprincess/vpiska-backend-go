package repository

import (
	"context"

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

type groupingResult struct {
	Name  string `bson:"_id"`
	Phone string `bson:"phone"`
}

func (r *userRepository) CheckNameAndPhone(ctx context.Context, name string, phone string) error {
	or := bson.D{{"$or", bson.A{
		bson.D{{"name", name}},
		bson.D{{"phone", phone}},
	}}}
	match := bson.D{{"$match", or}}
	group := bson.D{{"$group", bson.D{{"_id", "$name"}, {"phone", bson.D{{"$last", "$phone"}}}}}}
	pipeline := bson.A{match, group}
	cursor, aggregateErr := r.db.Aggregate(ctx, pipeline)

	if aggregateErr != nil {
		return aggregateErr
	}

	var result []groupingResult
	cursorErr := cursor.All(ctx, &result)

	if cursorErr != nil {
		return cursorErr
	}

	if result == nil {
		return nil
	}

	if len(result) == 2 {
		return domain.ErrNameAndPhoneAlreadyUse
	}

	if result[0].Name == name && result[0].Phone == phone {
		return domain.ErrNameAndPhoneAlreadyUse
	}

	if result[0].Phone == phone {
		return domain.ErrPhoneAlreadyUse
	}

	return domain.ErrNameAlreadyUse
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
	find := bson.D{{"_id", id}}
	return getUserByFilter(ctx, r.db, find)
}

func (r *userRepository) GetUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	find := bson.D{{"phone", phone}}
	return getUserByFilter(ctx, r.db, find)
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

func getUserByFilter(ctx context.Context, db *mongo.Collection, filter bson.D) (domain.User, error) {
	model := domain.User{}

	if err := db.FindOne(ctx, filter).Decode(&model); err != nil {
		if err == mongo.ErrNoDocuments {
			return model, domain.ErrUserNotFound
		}
		return model, err
	}

	return model, nil
}
