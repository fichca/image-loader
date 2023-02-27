package service

import (
	"context"
	"github.com/fichca/image-loader/internal/config"
	"github.com/minio/minio-go/v7"
	"io"
)

type fileService struct {
	minioCfg    *config.Minio
	minioClient *minio.Client
}

func NewFileService(minio *minio.Client, minioCfg *config.Minio) *fileService {
	return &fileService{
		minioClient: minio,
		minioCfg:    minioCfg,
	}
}

func (fs *fileService) AddFile(ctx context.Context, filename string, file io.Reader) error {
	_, err := fs.minioClient.PutObject(ctx, fs.minioCfg.Bucket, filename, file, -1, minio.PutObjectOptions{})

	return err
}
