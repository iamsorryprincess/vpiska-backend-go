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

func (s *mediaService) GetMetadata(ctx context.Context, id string) (FileMetadata, error) {
	media, err := s.repository.GetMedia(ctx, id)

	if err != nil {
		return FileMetadata{}, err
	}

	metadata := FileMetadata{
		ID:               media.ID,
		Name:             media.Name,
		ContentType:      media.ContentType,
		Size:             media.Size,
		LastModifiedDate: media.LastModifiedDate,
	}

	return metadata, nil
}

func (s *mediaService) GetFile(ctx context.Context, id string) (*FileData, error) {
	metadata, err := s.repository.GetMedia(ctx, id)
	filename := path + "/" + id

	if err != nil {
		if err == domain.ErrMediaNotFound {
			fileErr := os.Remove(filename)

			if os.IsNotExist(fileErr) {
				return nil, err
			}

			return nil, fileErr
		}

		return nil, err
	}

	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)

	if err != nil {
		if os.IsNotExist(err) {
			repoErr := s.repository.DeleteMedia(ctx, id)

			if repoErr != nil {
				return nil, repoErr
			}

			return nil, domain.ErrMediaNotFound
		}

		return nil, err
	}

	result := &FileData{
		ContentType: metadata.ContentType,
		Size:        metadata.Size,
		File:        file,
	}

	return result, nil
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
