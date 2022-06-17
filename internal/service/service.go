package service

import (
	"context"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/storage"
)

type CreateMediaInput struct {
	Name        string
	ContentType string
	Size        int64
	Data        []byte
}

type Media interface {
	Create(ctx context.Context, input CreateMediaInput) (string, error)
	Update(ctx context.Context, mediaId string, input CreateMediaInput) error
	GetMetadata(ctx context.Context, id string) (domain.Media, error)
	GetFile(ctx context.Context, id string) (domain.FileData, error)
	Delete(ctx context.Context, id string) error
}

type CreateUserInput struct {
	Name     string
	Phone    string
	Password string
}

type LoginUserInput struct {
	Phone    string
	Password string
}

type ChangePasswordInput struct {
	ID       string
	Password string
}

type UpdateUserInput struct {
	ID    string
	Name  string
	Phone string
}

type SetUserImageInput struct {
	UserID      string
	FileName    string
	ContentType string
	Size        int64
	FileData    []byte
}

type Users interface {
	Create(ctx context.Context, input CreateUserInput) (domain.UserLogin, error)
	Login(ctx context.Context, input LoginUserInput) (domain.UserLogin, error)
	Update(ctx context.Context, input UpdateUserInput) (string, error)
	ChangePassword(ctx context.Context, input ChangePasswordInput) (string, error)
	SetUserImage(ctx context.Context, input SetUserImageInput) (imageId string, accessToken string, err error)
}

type CreateEventInput struct {
	OwnerID     string
	Name        string
	Address     string
	Coordinates domain.Coordinates
}

type GetByRangeInput struct {
	HorizontalRange float64
	VerticalRange   float64
	Coordinates     domain.Coordinates
}
type AddUserInfoInput struct {
	EventID string
	UserID  string
}

type ChatMessageInput struct {
	EventID     string
	UserID      string
	UserName    string
	UserImageID string
	Message     string
}

type AddMediaInput struct {
	EventID     string
	UserID      string
	FileName    string
	ContentType string
	FileSize    int64
	FileData    []byte
}

type UpdateEventInput struct {
	UserID      string
	EventID     string
	Address     string
	Coordinates domain.Coordinates
}

type RemoveMediaInput struct {
	EventID string
	UserID  string
	MediaID string
}

type Events interface {
	Create(ctx context.Context, input CreateEventInput) (domain.EventInfo, error)
	Close(ctx context.Context, eventId string, userId string) error
	Update(ctx context.Context, input UpdateEventInput) error
	GetByID(ctx context.Context, id string) (domain.EventInfo, error)
	GetByRange(ctx context.Context, input GetByRangeInput) ([]domain.EventRangeData, error)
	AddUserInfo(ctx context.Context, input AddUserInfoInput) error
	RemoveUserInfo(ctx context.Context, eventId string, userId string) error
	SendChatMessage(ctx context.Context, input ChatMessageInput) error
	AddMedia(ctx context.Context, input AddMediaInput) error
	RemoveMedia(ctx context.Context, input RemoveMediaInput) error
}

type Subscriber interface {
	OnReceive(message []byte)
	OnClose()
}

type Publisher interface {
	Subscribe(eventId string, subscriber Subscriber)
	Unsubscribe(eventId string, subscriber Subscriber)
	Publish(eventId string, message []byte)
	Close(eventId string)
	CloseAll()
}

type Services struct {
	Media     Media
	Users     Users
	Events    Events
	Publisher Publisher
}

func NewServices(
	logger logger.Logger,
	repositories *repository.Repositories,
	hashManager hash.PasswordHashManager,
	auth auth.TokenManager,
	storage storage.FileStorage) (*Services, error) {
	media := newMediaService(repositories.Media, storage)
	pub := newPublisher()

	return &Services{
		Media:     media,
		Users:     newUserService(repositories.Users, repositories.Events, hashManager, auth, media),
		Events:    NewEventService(logger, repositories.Events, pub, media),
		Publisher: pub,
	}, nil
}
