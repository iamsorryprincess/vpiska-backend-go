package repository

import (
	"context"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository/mongo"
)

type Media interface {
	GetMedia(ctx context.Context, id string) (domain.Media, error)
	CreateMedia(ctx context.Context, media domain.Media) (string, error)
	UpdateMedia(ctx context.Context, media domain.Media) error
	DeleteMedia(ctx context.Context, id string) error
}

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

type Events interface {
	CreateEvent(ctx context.Context, event domain.Event) (string, error)
	GetEventById(ctx context.Context, id string) (domain.Event, error)
	GetEventByOwnerId(ctx context.Context, ownerId string) (domain.Event, error)
	GetEventsByRange(ctx context.Context, xLeft float64, xRight float64, yLeft float64, yRight float64) ([]domain.EventRangeData, error)
	UpdateEvent(ctx context.Context, id string, address string, coordinates domain.Coordinates) error
	RemoveEvent(ctx context.Context, id string) error
	AddMedia(ctx context.Context, id string, mediaInfo domain.MediaInfo) error
	RemoveMedia(ctx context.Context, eventId string, mediaId string) error
	AddUserInfo(ctx context.Context, eventId string, userInfo domain.UserInfo) error
	RemoveUserInfo(ctx context.Context, eventId string, userId string) error
	AddChatMessage(ctx context.Context, id string, chatMessage domain.ChatMessage) error
}

type Repositories struct {
	Media  Media
	Users  Users
	Events Events
}

type TestsCleaner interface {
	Clean() error
}

func NewRepositories(connectionString string, dbName string) (*Repositories, TestsCleaner, error) {
	media, users, events, cleaner, err := mongo.NewRepositories(connectionString, dbName)
	if err != nil {
		return nil, nil, err
	}

	return &Repositories{
		Media:  media,
		Users:  users,
		Events: events,
	}, cleaner, nil
}
