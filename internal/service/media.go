package service

import (
	"context"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/storage"
)

type mediaService struct {
	repository repository.Media
	storage    storage.FileStorage
}

func newMediaService(repository repository.Media, fileStorage storage.FileStorage) Media {
	return &mediaService{
		repository: repository,
		storage:    fileStorage,
	}
}

func (s *mediaService) Create(ctx context.Context, input CreateMediaInput) (string, error) {
	media := domain.Media{
		Name:             input.Name,
		ContentType:      input.ContentType,
		Size:             input.Size,
		LastModifiedDate: time.Now(),
	}

	mediaId, err := s.repository.CreateMedia(ctx, media)

	if err != nil {
		return "", err
	}

	err = s.storage.Upload(mediaId, input.Data)

	if err != nil {
		_ = s.repository.DeleteMedia(ctx, mediaId)
		return "", err
	}

	return mediaId, nil
}

func (s *mediaService) Update(ctx context.Context, mediaId string, input CreateMediaInput) error {
	err := s.repository.UpdateMedia(ctx, domain.Media{
		ID:               mediaId,
		Name:             input.Name,
		ContentType:      input.ContentType,
		Size:             input.Size,
		LastModifiedDate: time.Now(),
	})

	if err != nil {
		return err
	}

	return s.storage.Upload(mediaId, input.Data)
}

func (s *mediaService) GetMetadata(ctx context.Context, id string) (domain.Media, error) {
	return s.repository.GetMedia(ctx, id)
}

func (s *mediaService) GetFile(ctx context.Context, id string) (domain.FileData, error) {
	metadata, err := s.repository.GetMedia(ctx, id)

	if err != nil {
		return domain.FileData{}, err
	}

	fileData, err := s.storage.Get(id)

	if err != nil {
		return domain.FileData{}, err
	}

	return domain.FileData{
		ContentType: metadata.ContentType,
		Size:        metadata.Size,
		Data:        fileData,
	}, nil
}

func (s *mediaService) Delete(ctx context.Context, id string) error {
	err := s.repository.DeleteMedia(ctx, id)

	if err != nil {
		return err
	}

	err = s.storage.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
