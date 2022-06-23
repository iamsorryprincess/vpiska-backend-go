package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	Media  *Media
	Users  *Users
	Events *Events
}

func NewRepositories(connectionString string, dbName string) (*Repositories, *TestsCleaner, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))

	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, nil, err
	}

	db := client.Database(dbName)

	return &Repositories{
		Media:  newMedia(db, "media"),
		Users:  newUsers(db, "users"),
		Events: newEvents(db, "events"),
	}, newTestsCleaner(db), nil
}
