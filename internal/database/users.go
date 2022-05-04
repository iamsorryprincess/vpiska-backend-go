package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	db *mongo.Collection
}

func NewUserRepository(connectionString string, dbName string, collectionName string) (user.Repository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if connectionErr := client.Connect(ctx); connectionErr != nil {
		return nil, connectionErr
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	return &userRepository{
		db: collection,
	}, nil
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
		return user.NameAndPhoneAlreadyUse
	}

	if result[0].Name == name && result[0].Phone == phone {
		return user.NameAndPhoneAlreadyUse
	}

	if result[0].Phone == phone {
		return user.PhoneAlreadyUse
	}

	return user.NameAlreadyUse
}

func (r *userRepository) CreateUser(ctx context.Context, user *user.User) error {
	user.ID = uuid.New().String()
	_, err := r.db.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	find := bson.D{{"_id", id}}
	return getUserByFilter(ctx, r.db, find)
}

func (r *userRepository) GetUserByPhone(ctx context.Context, phone string) (*user.User, error) {
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
		return user.NotFound
	}

	return nil
}

func getUserByFilter(ctx context.Context, db *mongo.Collection, filter bson.D) (*user.User, error) {
	model := &user.User{}

	if err := db.FindOne(ctx, filter).Decode(&model); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.NotFound
		}
		return nil, err
	}

	return model, nil
}
