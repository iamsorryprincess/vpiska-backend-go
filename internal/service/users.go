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
	fileStorage Media
}

func newUserService(
	repository repository.Users,
	hashManager hash.PasswordHashManager,
	auth auth.TokenManager,
	fileStorage Media) Users {
	return &userService{
		repository:  repository,
		hashManager: hashManager,
		auth:        auth,
		fileStorage: fileStorage,
	}
}

func (s *userService) Create(ctx context.Context, input CreateUserInput) (LoginResponse, error) {
	namesCount, err := s.repository.GetNamesCount(ctx, input.Name)

	if err != nil {
		return LoginResponse{}, err
	}

	phonesCount, err := s.repository.GetPhonesCount(ctx, input.Phone)

	if err != nil {
		return LoginResponse{}, err
	}

	if namesCount > 0 && phonesCount > 0 {
		return LoginResponse{}, domain.ErrNameAndPhoneAlreadyUse
	}

	if namesCount > 0 {
		return LoginResponse{}, domain.ErrNameAlreadyUse
	}

	if phonesCount > 0 {
		return LoginResponse{}, domain.ErrPhoneAlreadyUse
	}

	model := domain.User{
		Name:      input.Name,
		PhoneCode: "+7",
		Phone:     input.Phone,
		Password:  s.hashManager.HashPassword(input.Password),
	}

	userId, err := s.repository.CreateUser(ctx, model)

	if err != nil {
		return LoginResponse{}, err
	}

	token, err := s.generateToken(userId, model.Name, model.ImageID)

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		ID:          userId,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		AccessToken: token,
	}, nil
}

func (s *userService) Login(ctx context.Context, input LoginUserInput) (LoginResponse, error) {
	model, err := s.repository.GetUserByPhone(ctx, input.Phone)

	if err != nil {
		return LoginResponse{}, err
	}

	if !s.hashManager.VerifyHashedPassword(model.Password, input.Password) {
		return LoginResponse{}, domain.ErrInvalidPassword
	}

	token, err := s.generateToken(model.ID, model.Name, model.ImageID)

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		ID:          model.ID,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		AccessToken: token,
	}, nil
}

func (s *userService) Update(ctx context.Context, input UpdateUserInput) error {
	if input.Name == "" && input.Phone == "" {
		return nil
	}

	if input.Name != "" && input.Phone != "" {
		namesCount, err := s.repository.GetNamesCount(ctx, input.Name)

		if err != nil {
			return err
		}

		phonesCount, err := s.repository.GetPhonesCount(ctx, input.Phone)

		if err != nil {
			return err
		}

		if namesCount > 0 && phonesCount > 0 {
			return domain.ErrNameAndPhoneAlreadyUse
		}

		if namesCount > 0 {
			return domain.ErrNameAlreadyUse
		}

		if phonesCount > 0 {
			return domain.ErrPhoneAlreadyUse
		}

		return s.repository.UpdateNameAndPhone(ctx, input.ID, input.Name, input.Phone)
	}

	if input.Name != "" {
		namesCount, err := s.repository.GetNamesCount(ctx, input.Name)

		if err != nil {
			return err
		}

		if namesCount > 0 {
			return domain.ErrNameAlreadyUse
		}

		return s.repository.UpdateName(ctx, input.ID, input.Name)
	}

	phonesCount, err := s.repository.GetPhonesCount(ctx, input.Phone)

	if err != nil {
		return err
	}

	if phonesCount > 0 {
		return domain.ErrPhoneAlreadyUse
	}

	return s.repository.UpdatePhone(ctx, input.ID, input.Phone)
}

func (s *userService) ChangePassword(ctx context.Context, input ChangePasswordInput) (LoginResponse, error) {
	model, err := s.repository.GetUserByID(ctx, input.ID)

	if err != nil {
		return LoginResponse{}, err
	}

	if err = s.repository.ChangePassword(ctx, input.ID, s.hashManager.HashPassword(input.Password)); err != nil {
		return LoginResponse{}, err
	}

	token, err := s.generateToken(model.ID, model.Name, model.ImageID)

	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		ID:          model.ID,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		AccessToken: token,
	}, nil
}

func (s *userService) SetUserImage(ctx context.Context, input *SetUserImageInput) (string, error) {
	user, err := s.repository.GetUserByID(ctx, input.UserID)

	if err != nil {
		return "", err
	}

	if user.ImageID == "" {
		imageId, err := s.fileStorage.Create(ctx, &CreateMediaInput{
			Name:        input.FileName,
			ContentType: input.ContentType,
			Size:        input.Size,
			Data:        input.FileData,
		})

		if err != nil {
			return "", err
		}

		err = s.repository.SetImageId(ctx, input.UserID, imageId)

		if err != nil {
			return "", err
		}

		return imageId, nil
	}

	err = s.fileStorage.Update(ctx, user.ImageID, &CreateMediaInput{
		Name:        input.FileName,
		ContentType: input.ContentType,
		Size:        input.Size,
		Data:        input.FileData,
	})

	if err != nil {
		return "", err
	}

	return user.ImageID, nil
}

func (s *userService) generateToken(id string, name string, imageId string) (string, error) {
	tokenInput := auth.TokenData{
		ID:      id,
		Name:    name,
		ImageID: imageId,
	}
	return s.auth.GetAccessToken(tokenInput)
}
