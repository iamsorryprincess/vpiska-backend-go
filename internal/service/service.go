package service

import (
	"context"
	"time"

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

type FileMetadata struct {
	ID               string
	Name             string
	ContentType      string
	Size             int64
	LastModifiedDate time.Time
}

type FileData struct {
	ContentType string
	Size        int64
	Data        []byte
}

type Media interface {
	Create(ctx context.Context, input *CreateMediaInput) (string, error)
	Update(ctx context.Context, mediaId string, input *CreateMediaInput) error
	GetMetadata(ctx context.Context, id string) (FileMetadata, error)
	GetFile(ctx context.Context, id string) (*FileData, error)
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
	Update(ctx context.Context, input UpdateUserInput) (domain.UserLogin, error)
	ChangePassword(ctx context.Context, input ChangePasswordInput) (domain.UserLogin, error)
	SetUserImage(ctx context.Context, input *SetUserImageInput) (domain.UserLogin, error)
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

type Events interface {
	Create(ctx context.Context, input CreateEventInput) (domain.Event, error)
	GetByID(ctx context.Context, id string) (domain.Event, error)
	GetByRange(ctx context.Context, input GetByRangeInput) ([]domain.EventRangeData, error)
}

type Services struct {
	Users  Users
	Media  Media
	Events Events
}

func NewServices(
	logger logger.Logger,
	repositories *repository.Repositories,
	hashManager hash.PasswordHashManager,
	auth auth.TokenManager,
	storage storage.FileStorage) (*Services, error) {
	media := newMediaService(repositories.Media, storage)

	return &Services{
		Users:  newUserService(repositories.Users, repositories.Events, hashManager, auth, media),
		Media:  media,
		Events: NewEventService(logger, repositories.Events),
	}, nil
}
