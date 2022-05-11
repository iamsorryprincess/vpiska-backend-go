package service

import (
	"context"

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

type Services struct {
	Users Users
}

func NewServices(
	repositories *repository.Repositories,
	hashManager hash.PasswordHashManager,
	auth auth.TokenManager) *Services {
	return &Services{
		Users: newUserService(repositories.Users, hashManager, auth),
	}
}
