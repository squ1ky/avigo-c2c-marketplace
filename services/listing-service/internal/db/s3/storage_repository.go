package s3

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
)

type StorageRepository interface {
	UploadFile(ctx context.Context, req UploadFileRequest) (*UploadFileResponse, error)
	GetPresignedDownloadURL(ctx context.Context, objectID string, expiryDuration time.Duration) (string, error)
	GetPresignedUploadURL(ctx context.Context, objectID string, expiryDuration time.Duration) (string, error)
	DeleteFile(ctx context.Context, objectID string) error
	DeleteFiles(ctx context.Context, objectIDs []string) error
	GetFileInfo(ctx context.Context, objectID string) (*FileInfo, error)
	CopyFile(ctx context.Context, srcObjectID, dstObjectID string) error
}

type UploadFileRequest struct {
	ObjectID    string
	Data        io.Reader
	Size        int64
	ContentType string
	Metadata    map[string]string
}

type UploadFileResponse struct {
	ObjectID    string
	Size        int64
	ContentType string
	ETag        string
	UploadedAt  time.Time
}

type FileInfo struct {
	ObjectID     string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
	Metadata     map[string]string
}

type MinIORepository struct {
	client *minio.Client
	config config.S3Config
}

func NewMinIORepository(cfg config.S3Config) (*MinIORepository, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{
			Region: cfg.Region,
		})

		if err != nil {
			return nil, err
		}
	}

	return &MinIORepository{
		client: client,
		config: cfg,
	}, nil
}

func (r *MinIORepository) UploadFile(ctx context.Context, objectKey string, data io.Reader, size int64, contentType string) error {
	_, err := r.client.PutObject(ctx, r.config.Bucket, objectKey, data, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (r *MinIORepository) GetPresignedURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	presignedURL, err := r.client.PresignedGetObject(ctx, r.config.Bucket, objectKey, expiry, nil)
	if err != nil {
		return "", err
	}

	url := presignedURL.String()

	if r.config.PublicURL != "" {
		url = strings.ReplaceAll(url, r.config.Endpoint, r.config.PublicURL)
	}

	return url, nil
}
