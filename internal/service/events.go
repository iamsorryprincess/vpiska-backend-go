package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

type eventService struct {
	logger     logger.Logger
	repository repository.Events
	publisher  Publisher
}

func NewEventService(logger logger.Logger, events repository.Events, publisher Publisher) Events {
	return &eventService{
		logger:     logger,
		repository: events,
		publisher:  publisher,
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

	s.publisher.Publish(input.EventID, []byte(fmt.Sprintf("usersCountUpdated/%d", len(event.Users)+1)))
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

			s.publisher.Publish(eventId, []byte(fmt.Sprintf("usersCountUpdated/%d", len(event.Users)-1)))
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
