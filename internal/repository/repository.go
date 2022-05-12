package repository

import (
	"context"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Users interface {
	CheckNameAndPhone(ctx context.Context, name string, phone string) error
	CreateUser(ctx context.Context, user domain.User) (string, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (domain.User, error)
	ChangePassword(ctx context.Context, id string, password string) error
}

type Media interface {
	CreateMedia(ctx context.Context, media domain.Media) (string, error)
	DeleteMedia(ctx context.Context, id string) error
}

type Repositories struct {
	Users Users
	Media Media
}

func NewRepositories(connectionString string, dbName string) (*Repositories, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &Repositories{
		Users: newMongoUsers(db, "users"),
		Media: newMongoMedia(db, "media"),
	}, nil
}
