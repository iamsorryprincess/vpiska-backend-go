package repository

import (
	"context"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Users interface {
	GetNamesCount(ctx context.Context, name string) (int64, error)
	GetPhonesCount(ctx context.Context, phone string) (int64, error)
	CreateUser(ctx context.Context, user domain.User) (string, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (domain.User, error)
	ChangePassword(ctx context.Context, id string, password string) error
	SetImageId(ctx context.Context, userId string, imageId string) error
	UpdateName(ctx context.Context, userId string, name string) error
	UpdatePhone(ctx context.Context, userId string, phone string) error
	UpdateNameAndPhone(ctx context.Context, userId string, name string, phone string) error
}

type Media interface {
	GetAll(ctx context.Context) ([]domain.Media, error)
	GetMedia(ctx context.Context, id string) (domain.Media, error)
	CreateMedia(ctx context.Context, media domain.Media) (string, error)
	UpdateMedia(ctx context.Context, media domain.Media) error
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
