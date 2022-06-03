package domain

import "time"

const EventStateOpened = 0
const EventStateClosed = 1

type EventState int

type Coordinates struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
}

type UserInfo struct {
	ID string `bson:"_id" json:"id"`
}

type MediaInfo struct {
	ID          string `bson:"_id"          json:"id"`
	ContentType string `bson:"content_type" json:"contentType"`
}

type ChatMessage struct {
	UserID      string `bson:"user_id"       json:"userId"`
	UserName    string `bson:"user_name"     json:"userName"`
	UserImageID string `bson:"user_image_id" json:"userImageId"`
	Message     string `bson:"message"       json:"message"`
}

type Event struct {
	ID           string        `bson:"_id"`
	OwnerID      string        `bson:"owner_id"`
	Name         string        `bson:"name"`
	Address      string        `bson:"address"`
	State        EventState    `bson:"state"`
	Coordinates  Coordinates   `bson:"coordinates"`
	CreatedAt    time.Time     `bson:"created_at"`
	Users        []UserInfo    `bson:"users"`
	Media        []MediaInfo   `bson:"media"`
	ChatMessages []ChatMessage `bson:"chat_messages"`
}

type EventRangeData struct {
	ID          string      `bson:"_id"         json:"id"`
	Name        string      `bson:"name"        json:"name"`
	UsersCount  int         `bson:"users_count" json:"usersCount"`
	Coordinates Coordinates `bson:"coordinates" json:"coordinates"`
}

type EventInfo struct {
	ID           string        `bson:"_id"           json:"id"`
	OwnerID      string        `bson:"owner_id"      json:"ownerId"`
	Name         string        `bson:"name"          json:"name"`
	Address      string        `bson:"address"       json:"address"`
	Coordinates  Coordinates   `bson:"coordinates"   json:"coordinates"`
	UsersCount   int           `bson:"users_count"   json:"usersCount"`
	Media        []MediaInfo   `bson:"media"         json:"media"`
	ChatMessages []ChatMessage `bson:"chat_messages" json:"chatMessages"`
}
