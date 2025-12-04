package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/db/postgres"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/db/s3"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"mime/multipart"
	"slices"
	"strings"
	"time"
)

type MediaService struct {
	storage    *s3.MediaStorage
	repository *postgres.MediaRepository
	config     config.S3Config
}

func NewMediaService(storage *s3.MediaStorage, repo *postgres.MediaRepository, cfg config.S3Config) *MediaService {
	return &MediaService{
		storage:    storage,
		repository: repo,
		config:     cfg,
	}
}

func (s *MediaService) UploadMedia(ctx context.Context, listingID uuid.UUID, file *multipart.FileHeader, order int) (*domain.ListingMedia, error) {
	contentType := file.Header.Get("Content-Type")
	if err := s.validateFile(file.Size, contentType); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	fileUrl, err := s.storage.UploadFile(ctx, s3.UploadRequest{
		ListingID:   listingID,
		Filename:    file.Filename,
		Data:        src,
		Size:        file.Size,
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to storage: %w", err)
	}

	media := &domain.ListingMedia{
		ListingID: listingID,
		FileURL:   fileUrl,
		FileType:  s.determineMediaType(contentType),
		MimeType:  contentType,
		Order:     order,
		CreatedAt: time.Now(),
	}

	if err := s.repository.Create(ctx, media); err != nil {
		_ = s.storage.DeleteFile(ctx, fileUrl)
		return nil, fmt.Errorf("failed to save media metadata: %w", err)
	}

	return media, nil
}

func (s *MediaService) DeleteMedia(ctx context.Context, mediaID uuid.UUID) error {
	media, err := s.repository.GetByID(ctx, mediaID)
	if err != nil {
		return fmt.Errorf("failed to get media: %w", err)
	}

	if err := s.storage.DeleteFile(ctx, media.FileURL); err != nil {
		return fmt.Errorf("failed to delete from storage: %w", err)
	}

	if err := s.repository.Delete(ctx, mediaID); err != nil {
		return fmt.Errorf("failed to delete media metadata: %w", err)
	}

	return nil
}

type MediaWithURL struct {
	domain.ListingMedia
	URL string
}

func (s *MediaService) GetMediaWithURLs(ctx context.Context, listingID uuid.UUID) ([]MediaWithURL, error) {
	media, err := s.repository.GetByListingID(ctx, listingID)
	if err != nil {
		return nil, err
	}

	result := make([]MediaWithURL, len(media))
	for i, m := range media {
		url, err := s.storage.GetPresignedDownloadURL(ctx, m.FileURL, 24*time.Hour)
		if err != nil {
			return nil, fmt.Errorf("failed to generate URL: %w", err)
		}

		result[i] = MediaWithURL{
			ListingMedia: m,
			URL:          url,
		}
	}

	return result, nil
}

func (s *MediaService) validateFile(size int64, contentType string) error {
	if size > s.config.MaxFileSize {
		return fmt.Errorf("file too large: max %d bytes (%.1fMB)", s.config.MaxFileSize, float64(s.config.MaxFileSize)/1024/1024)
	}

	if !slices.Contains(s.config.AllowedTypes, contentType) {
		return fmt.Errorf("unsupported file type: %s. Allowed: %v", contentType, s.config.AllowedTypes)
	}

	return nil
}

func (s *MediaService) determineMediaType(mimeType string) domain.MediaType {
	if strings.HasPrefix(mimeType, "image/") {
		return domain.MediaTypeImage
	}
	return domain.MediaTypeVideo
}
