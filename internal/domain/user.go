package domain

type User struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	PhoneCode string `bson:"phone_code"`
	Phone     string `bson:"phone"`
	ImageID   string `bson:"image_id"`
	Password  string `bson:"password"`
}

type UserLogin struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	ImageID     *string `json:"imageId"`
	EventID     *string `json:"eventId"`
	AccessToken string  `json:"accessToken"`
}
