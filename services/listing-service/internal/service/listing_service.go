package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/mapper"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/mongo"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/postgres"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/s3"
	"sync"
	"time"
)

type ListingService struct {
	listingRepo *pgrepo.ListingRepository
	mediaRepo   *pgrepo.MediaRepository
	charsRepo   *mngrepo.CharacteristicsRepository
	storage     *s3.MediaStorage
	txManager   pgrepo.TransactionManager
}

func NewListingService(
	listingRepo *pgrepo.ListingRepository,
	mediaRepo *pgrepo.MediaRepository,
	charsRepo *mngrepo.CharacteristicsRepository,
	storage *s3.MediaStorage,
	txManager pgrepo.TransactionManager,
) *ListingService {
	return &ListingService{
		listingRepo: listingRepo,
		mediaRepo:   mediaRepo,
		charsRepo:   charsRepo,
		storage:     storage,
		txManager:   txManager,
	}
}

func (s *ListingService) Create(ctx context.Context, input dto.CreateListingInput) (*dto.ListingResponse, error) {
	listingID := uuid.New()
	now := time.Now()

	listing := mapper.ToDomainListing(input, listingID, now)
	chars := mapper.ToDomainCharacteristics(input, listingID, now)

	mediaList := make([]domain.ListingMedia, len(input.Files))
	uploadedKeys := make([]string, 0, len(input.Files))

	errCh := make(chan error, len(input.Files))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, file := range input.Files {
		wg.Add(1)

		go func(idx int, f dto.FileInput) {
			defer wg.Done()

			objectKey, err := s.storage.UploadFile(ctx, s3.UploadRequest{
				ListingID:   listingID,
				Filename:    file.Name,
				Data:        file.Data,
				Size:        file.Size,
				ContentType: file.ContentType,
			})
			if err != nil {
				errCh <- fmt.Errorf("file %s upload failed: %w", f.Name, err)
				return
			}

			mu.Lock()
			uploadedKeys = append(uploadedKeys, objectKey)
			mu.Unlock()

			mediaList[idx] = domain.ListingMedia{
				ID:        uuid.New(),
				ListingID: listingID,
				FileURL:   objectKey,
				FileType:  domain.DetectMediaType(f.ContentType),
				MimeType:  f.ContentType,
				Order:     idx,
				CreatedAt: now,
			}
		}(i, file)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		firstErr := <-errCh
		go s.storage.DeleteFiles(context.Background(), uploadedKeys)
		return nil, firstErr
	}

	if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.listingRepo.Create(txCtx, listing); err != nil {
			return err
		}

		for _, media := range mediaList {
			if err := s.mediaRepo.Create(txCtx, &media); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		go s.storage.DeleteFiles(context.Background(), uploadedKeys)
		return nil, fmt.Errorf("failed to create listing (pg): %w", err)
	}

	if err := s.charsRepo.Create(ctx, chars); err != nil {
		_ = s.listingRepo.Delete(ctx, listingID)
		return nil, fmt.Errorf("failed to create listing (mongo): %w", err)
	}

	return mapper.ToListingResponse(listing, chars, nil), nil
}

func (s *ListingService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ListingResponse, error) {
	listing, err := s.listingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	media, _ := s.mediaRepo.GetByListingID(ctx, id)
	chars, _ := s.charsRepo.GetByListingID(ctx, id)

	go func() {
		_ = s.listingRepo.IncrementViews(context.Background(), id)
	}()

	return mapper.ToListingResponse(listing, chars, media), nil
}
