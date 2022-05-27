package service

import (
	"context"
	"errors"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
)

type userService struct {
	repository      repository.Users
	eventRepository repository.Events
	hashManager     hash.PasswordHashManager
	auth            auth.TokenManager
	fileStorage     Media
}

func newUserService(
	repository repository.Users,
	eventRepository repository.Events,
	hashManager hash.PasswordHashManager,
	auth auth.TokenManager,
	fileStorage Media) Users {
	return &userService{
		repository:      repository,
		eventRepository: eventRepository,
		hashManager:     hashManager,
		auth:            auth,
		fileStorage:     fileStorage,
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

	hashedPassword, err := s.hashManager.HashPassword(input.Password)

	if err != nil {
		return LoginResponse{}, err
	}

	model := domain.User{
		Name:      input.Name,
		PhoneCode: "+7",
		Phone:     input.Phone,
		Password:  hashedPassword,
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

	isPasswordMatched, err := s.hashManager.VerifyPassword(input.Password, model.Password)

	if err != nil {
		return LoginResponse{}, err
	}

	if !isPasswordMatched {
		return LoginResponse{}, domain.ErrInvalidPassword
	}

	token, err := s.generateToken(model.ID, model.Name, model.ImageID)

	if err != nil {
		return LoginResponse{}, err
	}

	event, err := s.eventRepository.GetEventByOwnerId(ctx, model.ID)

	if err != nil && !errors.Is(err, domain.ErrEventNotFound) {
		return LoginResponse{}, err
	}

	return LoginResponse{
		ID:          model.ID,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		EventID:     event.ID,
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

	hashedPassword, err := s.hashManager.HashPassword(model.Password)

	if err != nil {
		return LoginResponse{}, err
	}

	if err = s.repository.ChangePassword(ctx, input.ID, hashedPassword); err != nil {
		return LoginResponse{}, err
	}

	token, err := s.generateToken(model.ID, model.Name, model.ImageID)

	if err != nil {
		return LoginResponse{}, err
	}

	event, err := s.eventRepository.GetEventByOwnerId(ctx, model.ID)

	if err != nil && !errors.Is(err, domain.ErrEventNotFound) {
		return LoginResponse{}, err
	}

	return LoginResponse{
		ID:          model.ID,
		Name:        model.Name,
		Phone:       model.Phone,
		ImageID:     model.ImageID,
		EventID:     event.ID,
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
