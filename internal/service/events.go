package service

import (
	"context"
	"errors"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

type eventService struct {
	logger     logger.Logger
	repository repository.Events
}

func NewEventService(logger logger.Logger, events repository.Events) Events {
	return &eventService{
		logger:     logger,
		repository: events,
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
	return s.repository.GetEventById(ctx, id)
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
