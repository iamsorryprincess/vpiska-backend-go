package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

type eventService struct {
	logger      logger.Logger
	repository  repository.Events
	publisher   Publisher
	fileStorage Media
}

func NewEventService(logger logger.Logger, events repository.Events, publisher Publisher, fileStorage Media) Events {
	return &eventService{
		logger:      logger,
		repository:  events,
		publisher:   publisher,
		fileStorage: fileStorage,
	}
}

func (s *eventService) Create(ctx context.Context, input CreateEventInput) (domain.EventInfo, error) {
	_, err := s.repository.GetEventByOwnerId(ctx, input.OwnerID)

	if err != nil && !errors.Is(err, domain.ErrEventNotFound) {
		return domain.EventInfo{}, err
	}

	if err == nil {
		return domain.EventInfo{}, domain.ErrOwnerAlreadyHasEvent
	}

	event := domain.Event{
		OwnerID:      input.OwnerID,
		Name:         input.Name,
		Address:      input.Address,
		Coordinates:  input.Coordinates,
		State:        domain.EventStateOpened,
		CreatedAt:    time.Now(),
		Users:        []domain.UserInfo{},
		Media:        []domain.MediaInfo{},
		ChatMessages: []domain.ChatMessage{},
	}

	id, err := s.repository.CreateEvent(ctx, event)

	if err != nil {
		return domain.EventInfo{}, err
	}

	return domain.EventInfo{
		ID:           id,
		OwnerID:      event.OwnerID,
		Name:         event.Name,
		Address:      event.Address,
		Coordinates:  event.Coordinates,
		UsersCount:   len(event.Users),
		Media:        event.Media,
		ChatMessages: event.ChatMessages,
	}, nil
}

func (s *eventService) Close(ctx context.Context, eventId string, userId string) error {
	event, err := s.repository.GetEventById(ctx, eventId)

	if err != nil {
		return err
	}

	if event.OwnerID != userId {
		return domain.ErrUserIsNotOwner
	}

	err = s.repository.RemoveEvent(ctx, eventId)

	if err != nil {
		return err
	}

	s.publisher.Publish(eventId, []byte("closeEvent/"))
	s.publisher.Close(eventId)
	return nil
}

type eventUpdated struct {
	EventID     string             `json:"eventId"`
	Name        string             `json:"name"`
	Address     string             `json:"address"`
	UsersCount  int                `json:"usersCount"`
	Coordinates domain.Coordinates `json:"coordinates"`
}

func (s *eventService) Update(ctx context.Context, input UpdateEventInput) error {
	event, err := s.repository.GetEventById(ctx, input.EventID)

	if err != nil {
		return err
	}

	if event.OwnerID != input.UserID {
		return domain.ErrUserIsNotOwner
	}

	err = s.repository.UpdateEvent(ctx, input.EventID, input.Address, input.Coordinates)

	if err != nil {
		return err
	}

	data, err := json.Marshal(eventUpdated{
		EventID:     event.ID,
		Name:        event.Name,
		Address:     input.Address,
		UsersCount:  len(event.Users),
		Coordinates: event.Coordinates,
	})

	if err != nil {
		return err
	}

	s.publisher.Publish(input.EventID, []byte("eventUpdated/"+string(data)))
	return nil
}

func (s *eventService) GetByID(ctx context.Context, id string) (domain.EventInfo, error) {
	event, err := s.repository.GetEventById(ctx, id)

	if err != nil {
		return domain.EventInfo{}, err
	}

	return domain.EventInfo{
		ID:           event.ID,
		OwnerID:      event.OwnerID,
		Name:         event.Name,
		Address:      event.Address,
		Coordinates:  event.Coordinates,
		UsersCount:   len(event.Users),
		Media:        event.Media,
		ChatMessages: event.ChatMessages,
	}, nil
}

func (s *eventService) GetByRange(ctx context.Context, input GetByRangeInput) ([]domain.EventRangeData, error) {
	halfHorizontalRange := input.HorizontalRange / 2
	halfVerticalRange := input.VerticalRange / 2
	xLeft := input.Coordinates.X - halfHorizontalRange
	xRight := input.Coordinates.X + halfHorizontalRange
	yLeft := input.Coordinates.Y - halfVerticalRange
	yRight := input.Coordinates.Y + halfVerticalRange
	result, err := s.repository.GetEventsByRange(ctx, xLeft, xRight, yLeft, yRight)

	if err != nil {
		return nil, err
	}

	return result, err
}

func (s *eventService) AddUserInfo(ctx context.Context, input AddUserInfoInput) error {
	event, err := s.repository.GetEventById(ctx, input.EventID)

	if err != nil {
		return err
	}

	for _, userInfo := range event.Users {
		if userInfo.ID == input.UserID {
			return domain.ErrUserAlreadyExist
		}
	}

	err = s.repository.AddUserInfo(ctx, input.EventID, domain.UserInfo{
		ID: input.UserID,
	})

	if err != nil {
		return err
	}

	data, err := json.Marshal(eventUpdated{
		EventID:     event.ID,
		Name:        event.Name,
		Address:     event.Address,
		UsersCount:  len(event.Users) + 1,
		Coordinates: event.Coordinates,
	})

	if err != nil {
		return err
	}

	s.publisher.Publish(input.EventID, []byte("eventUpdated/"+string(data)))
	return nil
}

func (s *eventService) RemoveUserInfo(ctx context.Context, eventId string, userId string) error {
	event, err := s.repository.GetEventById(ctx, eventId)

	if err != nil {
		return err
	}

	for _, userInfo := range event.Users {
		if userInfo.ID == userId {
			err = s.repository.RemoveUserInfo(ctx, eventId, userId)

			if err != nil {
				return err
			}

			data, err := json.Marshal(eventUpdated{
				EventID:     event.ID,
				Name:        event.Name,
				Address:     event.Address,
				UsersCount:  len(event.Users) - 1,
				Coordinates: event.Coordinates,
			})

			if err != nil {
				return err
			}

			s.publisher.Publish(eventId, []byte("eventUpdated/"+string(data)))
			return nil
		}
	}

	return domain.ErrUserNotFound
}

func (s *eventService) SendChatMessage(ctx context.Context, input ChatMessageInput) error {
	chatMessage := domain.ChatMessage{
		UserID:      input.UserID,
		UserName:    input.UserName,
		UserImageID: input.UserImageID,
		Message:     input.Message,
	}
	err := s.repository.AddChatMessage(ctx, input.EventID, chatMessage)

	if err != nil {
		return err
	}

	data, err := json.Marshal(chatMessage)

	if err != nil {
		return err
	}

	s.publisher.Publish(input.EventID, []byte("chatMessage/"+string(data)))
	return nil
}

func (s *eventService) AddMedia(ctx context.Context, input AddMediaInput) error {
	event, err := s.repository.GetEventById(ctx, input.EventID)

	if err != nil {
		return err
	}

	if event.OwnerID != input.UserID {
		return domain.ErrUserIsNotOwner
	}

	mediaId, err := s.fileStorage.Create(ctx, CreateMediaInput{
		Name:        input.FileName,
		ContentType: input.ContentType,
		Size:        input.FileSize,
		Data:        input.FileData,
	})

	if err != nil {
		return err
	}

	mediaInfo := domain.MediaInfo{
		ID:          mediaId,
		ContentType: input.ContentType,
	}
	err = s.repository.AddMedia(ctx, input.EventID, mediaInfo)

	if err != nil {
		return err
	}

	data, err := json.Marshal(mediaInfo)

	if err != nil {
		return err
	}

	s.publisher.Publish(input.EventID, []byte("mediaAdded/"+string(data)))
	return nil
}

func (s *eventService) RemoveMedia(ctx context.Context, input RemoveMediaInput) error {
	event, err := s.repository.GetEventById(ctx, input.EventID)

	if err != nil {
		return err
	}

	if event.OwnerID != input.UserID {
		return domain.ErrUserIsNotOwner
	}

	for _, mediaInfo := range event.Media {
		if mediaInfo.ID == input.MediaID {
			err = s.fileStorage.Delete(ctx, input.MediaID)

			if err != nil {
				return err
			}

			err = s.repository.RemoveMedia(ctx, input.EventID, input.MediaID)

			if err != nil {
				return err
			}

			s.publisher.Publish(input.EventID, []byte("mediaRemoved/"+input.MediaID))
			return nil
		}
	}

	return domain.ErrMediaNotFound
}
