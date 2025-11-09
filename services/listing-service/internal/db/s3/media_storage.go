package s3

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
	"io"
	"strings"
	"time"
)

type MediaStorage struct {
	client *minio.Client
	config config.S3Config
}

func NewMediaStorage(cfg config.S3Config) (*MediaStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create bucket if not exists
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{
			Region: cfg.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &MediaStorage{
		client: client,
		config: cfg,
	}, nil
}

type FileInfo struct {
	ObjectKey    string
	Size         int64
	ContentType  string
	LastModified time.Time
	ETag         string
	Metadata     map[string]string
}

func (s *MediaStorage) UploadFile(ctx context.Context, listingID uuid.UUID, filename string, data io.Reader, size int64, contentType string) (string, error) {
	// Generate unique path: listings/{listing_id}/{uuid}_{filename}
	objectKey := fmt.Sprintf("listings/%s/%s_%s", listingID, uuid.New(), filename)

	_, err := s.client.PutObject(ctx, s.config.Bucket, objectKey, data, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return objectKey, nil
}

func (s *MediaStorage) GetPresignedDownloadURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	presignedURL, err := s.client.PresignedGetObject(ctx, s.config.Bucket, objectKey, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	url := presignedURL.String()

	// Replace internal endpoint on public (for docker/k8)
	if s.config.PublicURL != "" {
		url = strings.ReplaceAll(url, s.config.Endpoint, s.config.PublicURL)
	}

	return url, nil
}

func (s *MediaStorage) GetPresignedUploadURL(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	presignedURL, err := s.client.PresignedPutObject(ctx, s.config.Bucket, objectKey, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	url := presignedURL.String()

	if s.config.PublicURL != "" {
		url = strings.ReplaceAll(url, s.config.Endpoint, s.config.PublicURL)
	}

	return url, nil
}

func (s *MediaStorage) DeleteFile(ctx context.Context, objectKey string) error {
	err := s.client.RemoveObject(ctx, s.config.Bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *MediaStorage) DeleteFiles(ctx context.Context, objectKeys []string) error {
	objectsCh := make(chan minio.ObjectInfo, len(objectKeys))

	go func() {
		defer close(objectsCh)
		for _, key := range objectKeys {
			objectsCh <- minio.ObjectInfo{Key: key}
		}
	}()

	errorsCh := s.client.RemoveObjects(ctx, s.config.Bucket, objectsCh, minio.RemoveObjectsOptions{})

	for err := range errorsCh {
		if err.Err != nil {
			return fmt.Errorf("failed to delete file %s: %w", err.ObjectName, err.Err)
		}
	}

	return nil
}

func (s *MediaStorage) DeleteListingFiles(ctx context.Context, listingID uuid.UUID) error {
	prefix := fmt.Sprintf("listings/%s/", listingID)

	objectsCh := s.client.ListObjects(ctx, s.config.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	objectsToDelete := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsToDelete)
		for object := range objectsCh {
			if object.Err != nil {
				continue
			}
			objectsToDelete <- object
		}
	}()

	errorsCh := s.client.RemoveObjects(ctx, s.config.Bucket, objectsToDelete, minio.RemoveObjectsOptions{})

	for err := range errorsCh {
		if err.Err != nil {
			return fmt.Errorf("failed to delete object %s: %w", err.ObjectName, err.Err)
		}
	}

	return nil
}

func (s *MediaStorage) GetFileInfo(ctx context.Context, objectKey string) (*FileInfo, error) {
	stat, err := s.client.StatObject(ctx, s.config.Bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &FileInfo{
		ObjectKey:    stat.Key,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
		LastModified: stat.LastModified,
		ETag:         stat.ETag,
		Metadata:     stat.UserMetadata,
	}, nil
}

func (s *MediaStorage) CopyFile(ctx context.Context, srcKey, dstKey string) error {
	src := minio.CopySrcOptions{
		Bucket: s.config.Bucket,
		Object: srcKey,
	}

	dst := minio.CopyDestOptions{
		Bucket: s.config.Bucket,
		Object: dstKey,
	}

	_, err := s.client.CopyObject(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
