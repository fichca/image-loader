package service

import (
	"context"
	"fmt"
	"github.com/fichca/image-loader/internal/dto"
	"github.com/fichca/image-loader/internal/entity"
	"github.com/gofrs/uuid"
	"io"
)

type imageStorage interface {
	PutObject(ctx context.Context, image string, data io.Reader) error
	GetImageUrls(ctx context.Context, imageNames []string) ([]string, error)
	GetObjects(ctx context.Context, imageNames []string) ([]io.Reader, error)
}

type imageRepository interface {
	Add(ctx context.Context, modelImage entity.Image) error
	GetAllByUserId(ctx context.Context, userID int) ([]entity.Image, error)
}
type FileService struct {
	fileStorage     imageStorage
	imageRepository imageRepository
}

func NewFileService(fileStorage imageStorage, imageRepository imageRepository) *FileService {
	return &FileService{
		fileStorage:     fileStorage,
		imageRepository: imageRepository,
	}
}

func (fs *FileService) AddImage(ctx context.Context, image dto.Image) error {
	imageName, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("failed to generate image name: %w", err)
	}
	image.Name = imageName.String() + image.Extension

	err = fs.imageRepository.Add(ctx, toImageEntity(image))
	if err != nil {
		return fmt.Errorf("failed to save image data to db: %w", err)
	}

	err = fs.fileStorage.PutObject(ctx, image.Name, image.Data)
	if err != nil {
		return fmt.Errorf("failed to put image to fileStore: %w", err)
	}

	return nil
}

func (fs FileService) GetImageUrlsByUserId(ctx context.Context, userId int) ([]string, error) {
	images, err := fs.imageRepository.GetAllByUserId(ctx, userId)
	if err != nil {
		return []string{}, fmt.Errorf("failed to get images by userID: %w", err)
	}
	names := make([]string, 0)
	for _, image := range images {
		names = append(names, image.Name)
	}
	urls, err := fs.fileStorage.GetImageUrls(ctx, names)
	if err != nil {
		return urls, fmt.Errorf("failed to get image urls: %w", err)
	}
	return urls, nil
}

func toImageEntity(image dto.Image) entity.Image {
	return entity.Image{
		ID:        image.ID,
		UserID:    image.UserID,
		Name:      image.Name,
		Extension: image.Extension,
	}
}
