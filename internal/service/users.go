package service

import (
	"context"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/security"
)

type userService struct {
	repository repository.Users
	security   security.PasswordManager
	auth       auth.TokenManager
}

func newUserService(
	repository repository.Users,
	security security.PasswordManager,
	auth auth.TokenManager) Users {
	return &userService{
		repository: repository,
		security:   security,
		auth:       auth,
	}
}

func (s *userService) Create(ctx context.Context, input CreateUserInput) (LoginResponse, error) {
	checkError := s.repository.CheckNameAndPhone(ctx, input.Name, input.Phone)

	if checkError != nil {
		return LoginResponse{}, domain.MapUserError(checkError)
	}

	model := domain.User{
		Name:      input.Name,
		PhoneCode: "+7",
		Phone:     input.Phone,
		Password:  s.security.HashPassword(input.Password),
	}

	userId, createError := s.repository.CreateUser(ctx, model)

	if createError != nil {
		return LoginResponse{}, domain.MapUserError(createError)
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

	if !s.security.VerifyHashedPassword(model.Password, input.Password) {
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
	model, findErr := s.repository.GetUserByID(ctx, input.ID)

	if findErr != nil {
		return LoginResponse{}, domain.MapUserError(findErr)
	}

	if updateErr := s.repository.ChangePassword(ctx, input.ID, s.security.HashPassword(input.Password)); updateErr != nil {
		return LoginResponse{}, domain.MapUserError(updateErr)
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
