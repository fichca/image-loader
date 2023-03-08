package service

import (
	"context"
	"io"
)

type TelegramService struct {
	is imageObjectService
}

type imageObjectService interface {
	GetImageObjectsByUserId(ctx context.Context, userId int) ([]io.Reader, error)
}

func NewTelegramService(imageService imageObjectService) *TelegramService {
	return &TelegramService{
		is: imageService,
	}
}

func (t *TelegramService) GetImageObjects(ctx context.Context, userId int) ([]io.Reader, error) {
	objects, err := t.is.GetImageObjectsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return objects, nil
}
