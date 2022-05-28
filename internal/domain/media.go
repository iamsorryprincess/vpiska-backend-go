package domain

import "time"

type FileData struct {
	ContentType string
	Size        int64
	Data        []byte
}

type Media struct {
	ID               string    `bson:"_id"`
	Name             string    `bson:"name"`
	ContentType      string    `bson:"content_type"`
	Size             int64     `bson:"size"`
	LastModifiedDate time.Time `bson:"last_modified_date"`
}
