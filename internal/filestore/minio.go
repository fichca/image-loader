package filestore

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"time"
)

type Minio struct {
	minio  *minio.Client
	bucket string
}

func NewMinio(minioClient *minio.Client, bucket string) *Minio {
	return &Minio{
		minio:  minioClient,
		bucket: bucket,
	}
}

func (m *Minio) PutObject(ctx context.Context, image string, data io.Reader) error {
	_, err := m.minio.PutObject(ctx, m.bucket, image, data, -1, minio.PutObjectOptions{})

	return err
}

func (m *Minio) GetImageUrls(ctx context.Context, imageNames []string) ([]string, error) {
	var urls []string

	for i := range imageNames {
		url, err := m.minio.PresignedGetObject(ctx, m.bucket, imageNames[i], time.Hour*24, nil)
		if err != nil {
			return nil, err
		}

		urls = append(urls, url.String())
	}

	return urls, nil
}

func (m *Minio) GetObjects(ctx context.Context, imageNames []string) ([]io.Reader, error) {
	var objects []io.Reader

	for i := range imageNames {
		object, err := m.minio.GetObject(ctx, m.bucket, imageNames[i], minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}

	return objects, nil
}
