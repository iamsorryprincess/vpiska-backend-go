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

func (s *eventService) Create(ctx context.Context, input CreateEventInput) (domain.Event, error) {
	_, err := s.repository.GetEventByOwnerId(ctx, input.OwnerID)

	if err != nil && !errors.Is(err, domain.ErrEventNotFound) {
		return domain.Event{}, err
	}

	if err == nil {
		return domain.Event{}, domain.ErrOwnerAlreadyHasEvent
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
		return domain.Event{}, err
	}

	event.ID = id
	return event, nil
}

func (s *eventService) GetByID(ctx context.Context, id string) (domain.Event, error) {
	return s.repository.GetEventById(ctx, id)
}
