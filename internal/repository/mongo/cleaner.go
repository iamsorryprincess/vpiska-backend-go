package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type TestsCleaner struct {
	db *mongo.Database
}

func newTestsCleaner(db *mongo.Database) *TestsCleaner {
	return &TestsCleaner{
		db: db,
	}
}

func (c *TestsCleaner) Clean() error {
	return c.db.Drop(context.Background())
}
