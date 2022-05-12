package service

import (
	"context"
	"os"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
)

const path = "media"

type mediaService struct {
	repository repository.Media
}

func newMediaService(repository repository.Media) (Media, error) {
	err := initFileDir()

	if err != nil {
		return nil, err
	}

	return &mediaService{
		repository: repository,
	}, nil
}

func (s *mediaService) Create(ctx context.Context, input *CreateMediaInput) (string, error) {
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

	file, err := os.OpenFile(path+"/"+mediaId, os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		_ = s.repository.DeleteMedia(ctx, mediaId)
		return "", err
	}

	defer file.Close()
	_, err = file.ReadFrom(input.File)

	if err != nil {
		_ = s.repository.DeleteMedia(ctx, mediaId)
		return "", err
	}

	return mediaId, nil
}

func initFileDir() error {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			mkdirErr := os.Mkdir(path, 0777)
			if mkdirErr != nil {
				return mkdirErr
			}
			return nil
		}

		return err
	}

	return nil
}
