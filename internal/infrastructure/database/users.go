package database

import (
	"context"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/errors"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/interfaces"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	db *mongo.Collection
}

func InitUserRepository(connectionString string, dbName string, collectionName string) (interfaces.Repository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connectionError := client.Connect(ctx)

	if connectionError != nil {
		return nil, connectionError
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	return &UserRepository{
		db: collection,
	}, nil
}

func (r *UserRepository) Insert(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID().Hex()
	_, err := r.db.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := r.db.FindOne(ctx, bson.D{{"_id", id}}).Decode(user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.UserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	user := &models.User{}
	err := r.db.FindOne(ctx, bson.D{{"phone", phone}}).Decode(user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.UserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) ChangePassword(ctx context.Context, id string, password string) error {
	return nil
}

func (r *UserRepository) Update(ctx context.Context, id string, name string, phone string, imageID string) error {
	return nil
}

type groupingResult struct {
	Name  string `bson:"_id"`
	Phone string `bson:"phone"`
}

func (r *UserRepository) CheckNameAndPhone(ctx context.Context, name string, phone string) error {
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
		return errors.NameAndPhoneAlreadyUse
	}

	if result[0].Name == name && result[0].Phone == phone {
		return errors.NameAndPhoneAlreadyUse
	}

	if result[0].Phone == phone {
		return errors.PhoneAlreadyUse
	}

	return errors.NameAlreadyUse
}
