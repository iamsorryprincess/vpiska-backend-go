package domain

type User struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	PhoneCode string `bson:"phone_code"`
	Phone     string `bson:"phone"`
	ImageID   string `bson:"image_id"`
	Password  string `bson:"password"`
}
