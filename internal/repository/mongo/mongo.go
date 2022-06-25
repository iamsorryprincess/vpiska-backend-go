package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewRepositories(connectionString string, dbName string) (*Media, *Users, *Events, *TestsCleaner, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))

	if err != nil {
		return nil, nil, nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, nil, nil, nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, nil, nil, nil, err
	}

	db := client.Database(dbName)

	return newMedia(db, "media"), newUsers(db, "users"), newEvents(db, "events"), newTestsCleaner(db), nil
}
