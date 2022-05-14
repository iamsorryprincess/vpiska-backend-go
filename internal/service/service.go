package service

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
)

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

type Users interface {
	Create(ctx context.Context, input CreateUserInput) (LoginResponse, error)
	Login(ctx context.Context, input LoginUserInput) (LoginResponse, error)
	ChangePassword(ctx context.Context, input ChangePasswordInput) (LoginResponse, error)
}

type CreateMediaInput struct {
	Name        string
	ContentType string
	Size        int64
	File        io.Reader
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
	File        *os.File
}

type Media interface {
	Create(ctx context.Context, input *CreateMediaInput) (string, error)
	GetMetadata(ctx context.Context, id string) (FileMetadata, error)
	GetFile(ctx context.Context, id string) (*FileData, error)
}

type Services struct {
	Users Users
	Media Media
}

func NewServices(
	repositories *repository.Repositories,
	hashManager hash.PasswordHashManager,
	auth auth.TokenManager) (*Services, error) {

	media, err := newMediaService(repositories.Media)

	if err != nil {
		return nil, err
	}

	return &Services{
		Users: newUserService(repositories.Users, hashManager, auth),
		Media: media,
	}, nil
}
