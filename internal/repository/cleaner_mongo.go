package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCleaner struct {
	db *mongo.Database
}

func newMongoTestsCleaner(db *mongo.Database) TestsCleaner {
	return &mongoCleaner{
		db: db,
	}
}

func (c *mongoCleaner) Clean() error {
	return c.db.Drop(context.Background())
}
