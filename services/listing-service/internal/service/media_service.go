package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	pgrepo "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/postgres"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/s3"
	"mime/multipart"
	"slices"
	"strings"
	"time"
)

type MediaService struct {
	storage    *s3.MediaStorage
	repository *pgrepo.MediaRepository
	config     config.S3Config
}

func NewMediaService(storage *s3.MediaStorage, repo *pgrepo.MediaRepository, cfg config.S3Config) *MediaService {
	return &MediaService{
		storage:    storage,
		repository: repo,
		config:     cfg,
	}
}

type UploadResult struct {
	ID  uuid.UUID `json:"id"`
	URL string    `json:"url"`
}

// UploadTempFile uploads a file into temp storage (S3 + media_uploads)
func (s *MediaService) UploadTempFile(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) (*UploadResult, error) {
	contentType := fileHeader.Header.Get("Content-Type")
	if err := s.validateFile(fileHeader.Size, contentType); err != nil {
		return nil, err
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	fileID := uuid.New()
	key := fmt.Sprintf("uploads/%s/%s_%s", userID.String(), fileID.String(), fileHeader.Filename)

	uploadedKey, err := s.storage.UploadRaw(ctx, key, src, fileHeader.Size, contentType)
	if err != nil {
		return nil, fmt.Errorf("s3 upload failed: %w", err)
	}

	tempMedia := &pgrepo.TempMedia{
		ID:        fileID,
		UserID:    userID,
		S3Key:     uploadedKey,
		Filename:  fileHeader.Filename,
		MimeType:  contentType,
		Size:      fileHeader.Size,
		CreatedAt: time.Now(),
	}
	if err := s.repository.CreateTemp(ctx, tempMedia); err != nil {
		return nil, fmt.Errorf("db save failed: %w", err)
	}

	url := uploadedKey
	if s.config.PublicURL != "" {
		url = fmt.Sprintf("%s/%s", strings.TrimRight(s.config.PublicURL, "/"), uploadedKey)
	}

	return &UploadResult{
		ID:  fileID,
		URL: url,
	}, nil
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
