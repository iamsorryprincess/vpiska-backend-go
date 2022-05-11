package service

import (
	"context"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
)

type userService struct {
	repository  repository.Users
	hashManager hash.PasswordHashManager
	auth        auth.TokenManager
}

func newUserService(
	repository repository.Users,
	hashManager hash.PasswordHashManager,
	auth auth.TokenManager) Users {
	return &userService{
		repository:  repository,
		hashManager: hashManager,
		auth:        auth,
	}
}

func (s *userService) Create(ctx context.Context, input CreateUserInput) (LoginResponse, error) {
	err := s.repository.CheckNameAndPhone(ctx, input.Name, input.Phone)

	if err != nil {
		return LoginResponse{}, domain.MapUserError(err)
	}

	model := domain.User{
		Name:      input.Name,
		PhoneCode: "+7",
		Phone:     input.Phone,
		Password:  s.hashManager.HashPassword(input.Password),
	}

	userId, err := s.repository.CreateUser(ctx, model)

	if err != nil {
		return LoginResponse{}, domain.MapUserError(err)
	}

	tokenInput := auth.CreateTokenInput{
		ID:      userId,
		Name:    model.Name,
		ImageID: model.ImageID,
	}

	return LoginResponse{
		ID:          userId,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		AccessToken: s.auth.GetAccessToken(tokenInput),
	}, nil
}

func (s *userService) Login(ctx context.Context, input LoginUserInput) (LoginResponse, error) {
	model, err := s.repository.GetUserByPhone(ctx, input.Phone)

	if err != nil {
		return LoginResponse{}, domain.MapUserError(err)
	}

	if !s.hashManager.VerifyHashedPassword(model.Password, input.Password) {
		return LoginResponse{}, domain.ErrInvalidPassword
	}

	tokenInput := auth.CreateTokenInput{
		ID:      model.ID,
		Name:    model.Name,
		ImageID: model.ImageID,
	}

	return LoginResponse{
		ID:          model.ID,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		AccessToken: s.auth.GetAccessToken(tokenInput),
	}, nil
}

func (s *userService) ChangePassword(ctx context.Context, input ChangePasswordInput) (LoginResponse, error) {
	model, err := s.repository.GetUserByID(ctx, input.ID)

	if err != nil {
		return LoginResponse{}, domain.MapUserError(err)
	}

	if err = s.repository.ChangePassword(ctx, input.ID, s.hashManager.HashPassword(input.Password)); err != nil {
		return LoginResponse{}, domain.MapUserError(err)
	}

	tokenInput := auth.CreateTokenInput{
		ID:      model.ID,
		Name:    model.Name,
		ImageID: model.ImageID,
	}

	return LoginResponse{
		ID:          model.ID,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		AccessToken: s.auth.GetAccessToken(tokenInput),
	}, nil
}
