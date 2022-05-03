package service

import (
	"context"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user"
)

type UserService struct {
	repository user.Repository
	security   user.SecurityProvider
	identity   user.IdentityProvider
}

func NewUserService(
	repository user.Repository,
	security user.SecurityProvider,
	identity user.IdentityProvider) *UserService {
	return &UserService{
		repository: repository,
		security:   security,
		identity:   identity,
	}
}

func (s *UserService) Create(ctx context.Context, name string, phone string, password string) (*user.LoginResponse, error) {
	checkError := s.repository.CheckNameAndPhone(ctx, name, phone)

	if checkError != nil {
		return nil, user.MapError(checkError)
	}

	model := &user.User{
		Name:      name,
		PhoneCode: "+7",
		Phone:     phone,
		Password:  s.security.HashPassword(password),
	}

	createError := s.repository.CreateUser(ctx, model)

	if createError != nil {
		return nil, user.MapError(createError)
	}

	return &user.LoginResponse{
		ID:          model.ID,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		AccessToken: s.identity.GetAccessToken(model),
	}, nil
}

func (s *UserService) Login(ctx context.Context, phone string, password string) (*user.LoginResponse, error) {
	model, err := s.repository.GetUserByPhone(ctx, phone)

	if err != nil {
		return nil, user.MapError(err)
	}

	if !s.security.VerifyHashedPassword(model.Password, password) {
		return nil, user.InvalidPassword
	}

	return &user.LoginResponse{
		ID:          model.ID,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		AccessToken: s.identity.GetAccessToken(model),
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id string, password string) (*user.LoginResponse, error) {
	model, findErr := s.repository.GetUserByID(ctx, id)

	if findErr != nil {
		return nil, user.MapError(findErr)
	}

	if updateErr := s.repository.ChangePassword(ctx, id, s.security.HashPassword(password)); updateErr != nil {
		return nil, user.MapError(updateErr)
	}

	return &user.LoginResponse{
		ID:          model.ID,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		AccessToken: s.identity.GetAccessToken(model),
	}, nil
}
