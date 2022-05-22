package domain

type Coordinates struct {
	X float64 `bson:"x"`
	Y float64 `bson:"y"`
}

type UserInfo struct {
	ID string `bson:"_id"`
}

type MediaInfo struct {
	ID          string `bson:"_id"`
	ContentType string `bson:"content_type"`
}

type ChatMessage struct {
	UserID      string `bson:"user_id"`
	UserName    string `bson:"user_name"`
	UserImageID string `bson:"user_image_id"`
	Message     string `bson:"message"`
}

type Event struct {
	ID           string        `bson:"_id"`
	OwnerID      string        `bson:"owner_id"`
	Name         string        `bson:"name"`
	Address      string        `bson:"address"`
	Coordinates  Coordinates   `bson:"coordinates"`
	Users        []UserInfo    `bson:"users"`
	Media        []MediaInfo   `bson:"media"`
	ChatMessages []ChatMessage `bson:"chat_messages"`
}
