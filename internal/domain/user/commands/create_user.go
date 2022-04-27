package commands

import (
	"context"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/interfaces"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/models"
)

type CreateUserCommand struct {
	Name            string
	Phone           string
	Password        string
	ConfirmPassword string
}

type CreateUserHandler struct {
	repository       interfaces.Repository
	passwordProvider interfaces.PasswordHashProvider
	tokenProvider    interfaces.IdentityProvider
}

func InitCreateUserHandler(
	repository interfaces.Repository,
	passwordProvider interfaces.PasswordHashProvider,
	tokenProvider interfaces.IdentityProvider) *CreateUserHandler {
	return &CreateUserHandler{
		repository:       repository,
		passwordProvider: passwordProvider,
		tokenProvider:    tokenProvider,
	}
}

func (h *CreateUserHandler) Handle(ctx context.Context, command *CreateUserCommand) (*models.UserResponse, error) {
	checkError := h.repository.CheckNameAndPhone(ctx, command.Name, command.Phone)

	if checkError != nil {
		return nil, checkError
	}

	user := &models.User{
		Name:      command.Name,
		PhoneCode: "+7",
		Phone:     command.Phone,
		Password:  h.passwordProvider.HashPassword(command.Password),
	}

	insertError := h.repository.Insert(ctx, user)

	if insertError != nil {
		return nil, insertError
	}

	return &models.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Phone:       user.Phone,
		ImageID:     user.ImageID,
		AccessToken: h.tokenProvider.GetAccessToken(user),
	}, nil
}
