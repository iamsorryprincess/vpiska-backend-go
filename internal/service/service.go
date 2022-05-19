package service

import (
	"context"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
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

type LoginResponse struct {
	ID          string
	Name        string
	Phone       string
	ImageID     string
	AccessToken string
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
	Create(ctx context.Context, input CreateUserInput) (LoginResponse, error)
	Login(ctx context.Context, input LoginUserInput) (LoginResponse, error)
	Update(ctx context.Context, input UpdateUserInput) error
	ChangePassword(ctx context.Context, input ChangePasswordInput) (LoginResponse, error)
	SetUserImage(ctx context.Context, input *SetUserImageInput) (string, error)
}

type Services struct {
	Users Users
	Media Media
}

func NewServices(
	repositories *repository.Repositories,
	hashManager hash.PasswordHashManager,
	auth auth.TokenManager,
	storage storage.FileStorage) (*Services, error) {
	media := newMediaService(repositories.Media, storage)

	return &Services{
		Users: newUserService(repositories.Users, hashManager, auth, media),
		Media: media,
	}, nil
}
