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

func (s *userService) Create(ctx context.Context, input CreateUserInput) (domain.UserLogin, error) {
	namesCount, err := s.repository.GetNamesCount(ctx, input.Name)

	if err != nil {
		return domain.UserLogin{}, err
	}

	phonesCount, err := s.repository.GetPhonesCount(ctx, input.Phone)

	if err != nil {
		return domain.UserLogin{}, err
	}

	if namesCount > 0 && phonesCount > 0 {
		return domain.UserLogin{}, domain.ErrNameAndPhoneAlreadyUse
	}

	if namesCount > 0 {
		return domain.UserLogin{}, domain.ErrNameAlreadyUse
	}

	if phonesCount > 0 {
		return domain.UserLogin{}, domain.ErrPhoneAlreadyUse
	}

	hashedPassword, err := s.hashManager.HashPassword(input.Password)

	if err != nil {
		return domain.UserLogin{}, err
	}

	model := domain.User{
		Name:      input.Name,
		PhoneCode: "+7",
		Phone:     input.Phone,
		Password:  hashedPassword,
	}

	userId, err := s.repository.CreateUser(ctx, model)

	if err != nil {
		return domain.UserLogin{}, err
	}

	return s.createLogin(userId, model.Name, model.Phone, model.ImageID, "")
}

func (s *userService) Login(ctx context.Context, input LoginUserInput) (domain.UserLogin, error) {
	model, err := s.repository.GetUserByPhone(ctx, input.Phone)

	if err != nil {
		return domain.UserLogin{}, err
	}

	isPasswordMatched, err := s.hashManager.VerifyPassword(input.Password, model.Password)

	if err != nil {
		return domain.UserLogin{}, err
	}

	if !isPasswordMatched {
		return domain.UserLogin{}, domain.ErrInvalidPassword
	}

	event, err := s.eventRepository.GetEventByOwnerId(ctx, model.ID)

	if err != nil && !errors.Is(err, domain.ErrEventNotFound) {
		return domain.UserLogin{}, err
	}

	return s.createLogin(model.ID, model.Name, model.Phone, model.ImageID, event.ID)
}

func (s *userService) Update(ctx context.Context, input UpdateUserInput) (string, error) {
	model, err := s.repository.GetUserByID(ctx, input.ID)

	if err != nil {
		return "", err
	}

	if input.Name == "" && input.Phone == "" {
		return s.auth.GetAccessToken(auth.TokenData{
			ID:      input.ID,
			Name:    input.Name,
			ImageID: model.ImageID,
		})
	}

	if input.Name != "" && input.Phone != "" {
		return s.updateNameAndPhone(ctx, input.ID, input.Name, input.Phone, model.ImageID)
	}

	if input.Name != "" {
		return s.updateName(ctx, input.ID, input.Name, model.ImageID)
	}

	return s.updatePhone(ctx, input.ID, model.Name, input.Phone, model.ImageID)
}

func (s *userService) ChangePassword(ctx context.Context, input ChangePasswordInput) (string, error) {
	model, err := s.repository.GetUserByID(ctx, input.ID)

	if err != nil {
		return "", err
	}

	hashedPassword, err := s.hashManager.HashPassword(model.Password)

	if err != nil {
		return "", err
	}

	if err = s.repository.ChangePassword(ctx, input.ID, hashedPassword); err != nil {
		return "", err
	}

	return s.auth.GetAccessToken(auth.TokenData{
		ID:      model.ID,
		Name:    model.Name,
		ImageID: model.ImageID,
	})
}

func (s *userService) SetUserImage(ctx context.Context, input *SetUserImageInput) (imageId string, accessToken string, err error) {
	user, err := s.repository.GetUserByID(ctx, input.UserID)

	if err != nil {
		return "", "", err
	}

	if user.ImageID == "" {
		imageId, err = s.fileStorage.Create(ctx, &CreateMediaInput{
			Name:        input.FileName,
			ContentType: input.ContentType,
			Size:        input.Size,
			Data:        input.FileData,
		})

		if err != nil {
			return "", "", err
		}

		err = s.repository.SetImageId(ctx, input.UserID, imageId)

		if err != nil {
			return "", "", err
		}

		accessToken, err = s.auth.GetAccessToken(auth.TokenData{
			ID:      input.UserID,
			Name:    user.Name,
			ImageID: imageId,
		})

		if err != nil {
			return "", "", err
		}

		return imageId, accessToken, nil
	}

	err = s.fileStorage.Update(ctx, user.ImageID, &CreateMediaInput{
		Name:        input.FileName,
		ContentType: input.ContentType,
		Size:        input.Size,
		Data:        input.FileData,
	})

	if err != nil {
		return "", "", err
	}

	accessToken, err = s.auth.GetAccessToken(auth.TokenData{
		ID:      input.UserID,
		Name:    user.Name,
		ImageID: imageId,
	})

	if err != nil {
		return "", "", err
	}

	return user.ImageID, accessToken, nil
}

func (s *userService) createLogin(id string, name string, phone string, imageId string, eventId string) (domain.UserLogin, error) {
	token, err := s.auth.GetAccessToken(auth.TokenData{
		ID:      id,
		Name:    name,
		ImageID: imageId,
	})

	if err != nil {
		return domain.UserLogin{}, err
	}

	result := domain.UserLogin{
		ID:          id,
		Name:        name,
		Phone:       phone,
		AccessToken: token,
	}

	if imageId != "" {
		result.ImageID = &imageId
	}

	if eventId != "" {
		result.EventID = &eventId
	}

	return result, nil
}

func (s *userService) updateNameAndPhone(ctx context.Context, userId string, name string, phone string, imageId string) (string, error) {
	namesCount, err := s.repository.GetNamesCount(ctx, name)

	if err != nil {
		return "", err
	}

	phonesCount, err := s.repository.GetPhonesCount(ctx, phone)

	if err != nil {
		return "", err
	}

	if namesCount > 0 && phonesCount > 0 {
		return "", domain.ErrNameAndPhoneAlreadyUse
	}

	if namesCount > 0 {
		return "", domain.ErrNameAlreadyUse
	}

	if phonesCount > 0 {
		return "", domain.ErrPhoneAlreadyUse
	}

	if err = s.repository.UpdateNameAndPhone(ctx, userId, name, phone); err != nil {
		return "", err
	}

	return s.auth.GetAccessToken(auth.TokenData{
		ID:      userId,
		Name:    name,
		ImageID: imageId,
	})
}

func (s *userService) updateName(ctx context.Context, userId string, name string, imageId string) (string, error) {
	namesCount, err := s.repository.GetNamesCount(ctx, name)

	if err != nil {
		return "", err
	}

	if namesCount > 0 {
		return "", err
	}

	if err = s.repository.UpdateName(ctx, userId, name); err != nil {
		return "", err
	}

	return s.auth.GetAccessToken(auth.TokenData{
		ID:      userId,
		Name:    name,
		ImageID: imageId,
	})
}

func (s *userService) updatePhone(ctx context.Context, userId string, name string, phone string, imageId string) (string, error) {
	phonesCount, err := s.repository.GetPhonesCount(ctx, phone)

	if err != nil {
		return "", err
	}

	if phonesCount > 0 {
		return "", err
	}

	if err = s.repository.UpdatePhone(ctx, userId, phone); err != nil {
		return "", err
	}

	return s.auth.GetAccessToken(auth.TokenData{
		ID:      userId,
		Name:    name,
		ImageID: imageId,
	})
}
