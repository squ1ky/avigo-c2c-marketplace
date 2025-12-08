package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	mngrepo "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/mongo"
	pgrepo "github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/postgres"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/s3"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/search"

	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/mapper"
)

type ListingService struct {
	listingRepo *pgrepo.ListingRepository
	mediaRepo   *pgrepo.MediaRepository
	charsRepo   *mngrepo.CharacteristicsRepository
	storage     *s3.MediaStorage
	search      *search.Client
	txManager   pgrepo.TransactionManager
}

func NewListingService(
	listingRepo *pgrepo.ListingRepository,
	mediaRepo *pgrepo.MediaRepository,
	charsRepo *mngrepo.CharacteristicsRepository,
	storage *s3.MediaStorage,
	search *search.Client,
	txManager pgrepo.TransactionManager,
) *ListingService {
	return &ListingService{
		listingRepo: listingRepo,
		mediaRepo:   mediaRepo,
		charsRepo:   charsRepo,
		storage:     storage,
		search:      search,
		txManager:   txManager,
	}
}

func (s *ListingService) Create(ctx context.Context, input dto.CreateListingInput) (*dto.ListingResponse, error) {
	listingID := uuid.New()
	now := time.Now()

	listing := mapper.ToDomainListing(input, listingID, now)
	chars := mapper.ToDomainCharacteristics(input, listingID, now)

	var mediaList []domain.ListingMedia
	if len(input.MediaIDs) > 0 {
		tempFiles, err := s.mediaRepo.GetTempByIDs(ctx, input.MediaIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetchtemp media: %w", err)
		}

		if len(tempFiles) != len(input.MediaIDs) {
			return nil, fmt.Errorf("invalid media_ids: some files not found or access denied")
		}

		for _, tf := range tempFiles {
			if tf.UserID != input.UserID {
				return nil, fmt.Errorf("access denied to media file %s", tf.ID)
			}
		}

		mediaList = make([]domain.ListingMedia, len(tempFiles))
		for i, tmp := range tempFiles {
			mediaList[i] = domain.ListingMedia{
				ID:        uuid.New(),
				ListingID: listingID,
				FileURL:   tmp.S3Key,
				FileType:  domain.DetectMediaType(tmp.MimeType),
				Order:     i,
				CreatedAt: now,
			}
		}
	}

	if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.listingRepo.Create(txCtx, listing); err != nil {
			return err
		}

		if len(mediaList) > 0 {
			if err := s.mediaRepo.BulkCreate(txCtx, mediaList); err != nil {
				return err
			}

			if err := s.mediaRepo.DeleteTemp(txCtx, input.MediaIDs); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to create listing (pg): %w", err)
	}

	if err := s.charsRepo.Create(ctx, chars); err != nil {
		_ = s.listingRepo.Delete(ctx, listingID)
		return nil, fmt.Errorf("failed to create listing (mongo): %w", err)
	}

	s.enqueueIndexing(ctx, listing, chars)

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

// Update verifies user permissions, then compiles a final list of media files
// by merging existing (already attached) and new (from temp storage) files, arranging them in the client-specified order.
func (s *ListingService) Update(ctx context.Context, input dto.UpdateListingInput) (*dto.ListingResponse, error) {
	existing, err := s.listingRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if existing.UserID != input.UserID {
		return nil, ErrForbidden
	}

	oldMediaList, err := s.mediaRepo.GetByListingID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	oldMediaMap := make(map[uuid.UUID]domain.ListingMedia)
	oldKeysMap := make(map[string]bool)
	for _, m := range oldMediaList {
		oldMediaMap[m.ID] = m
		oldKeysMap[m.FileURL] = true
	}

	var tempIDs []uuid.UUID
	for _, id := range input.MediaIDs {
		if _, exists := oldMediaMap[id]; !exists {
			tempIDs = append(tempIDs, id)
		}
	}

	tempMediaMap := make(map[uuid.UUID]pgrepo.TempMedia)
	if len(tempIDs) > 0 {
		temps, err := s.mediaRepo.GetTempByIDs(ctx, tempIDs)
		if err != nil {
			return nil, err
		}
		for _, t := range temps {
			tempMediaMap[t.ID] = t
		}
	}

	finalMediaList := make([]domain.ListingMedia, 0, len(input.MediaIDs))
	keptKeys := make(map[string]bool)

	for i, id := range input.MediaIDs {
		var media domain.ListingMedia

		if old, exists := oldMediaMap[id]; exists {
			media = old
			media.ID = uuid.New()
		} else if temp, exists := tempMediaMap[id]; exists {
			if temp.UserID != input.UserID {
				return nil, ErrForbidden
			}
			media = domain.ListingMedia{
				ID:        uuid.New(),
				ListingID: input.ID,
				FileURL:   temp.S3Key,
				FileType:  domain.DetectMediaType(temp.MimeType),
				MimeType:  temp.MimeType,
				CreatedAt: time.Now(),
			}
		} else {
			continue
		}

		media.Order = i
		finalMediaList = append(finalMediaList, media)
		keptKeys[media.FileURL] = true
	}

	existing.Title = input.Title
	existing.Description = input.Description
	existing.Price = input.Price
	existing.Currency = input.Currency
	existing.CategoryID = input.CategoryID
	existing.UpdatedAt = time.Now()

	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.listingRepo.Update(txCtx, existing); err != nil {
			return err
		}

		if err := s.mediaRepo.DeleteAllByListingID(txCtx, input.ID); err != nil {
			return err
		}

		if len(finalMediaList) > 0 {
			if err := s.mediaRepo.BulkCreate(txCtx, finalMediaList); err != nil {
				return err
			}
		}

		if len(tempIDs) > 0 {
			if err := s.mediaRepo.DeleteTemp(txCtx, tempIDs); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	var keysToDelete []string
	for key := range oldKeysMap {
		if !keptKeys[key] {
			keysToDelete = append(keysToDelete, key)
		}
	}

	if len(keysToDelete) > 0 {
		go s.storage.DeleteFiles(context.Background(), keysToDelete)
	}

	chars, _ := s.charsRepo.GetByListingID(ctx, existing.ID)
	s.enqueueIndexing(ctx, existing, chars)

	return mapper.ToListingResponse(existing, nil, finalMediaList), nil
}

func (s *ListingService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	listing, err := s.listingRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if listing.UserID != userID {
		return ErrForbidden
	}

	mediaList, err := s.mediaRepo.GetByListingID(ctx, id)
	if err != nil {
		return err
	}

	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.mediaRepo.DeleteAllByListingID(txCtx, id); err != nil {
			return err
		}
		if err := s.listingRepo.Delete(txCtx, id); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	_ = s.charsRepo.Delete(ctx, id)

	if len(mediaList) > 0 {
		keys := make([]string, len(mediaList))
		for i, m := range mediaList {
			keys[i] = m.FileURL
		}

		go func() {
			if err := s.storage.DeleteFiles(context.Background(), keys); err != nil {
				slog.Warn("failed to cleanup s3 files: %v", err)
			}
		}()
	}

	s.enqueueDeletion(ctx, id)

	return nil
}

func (s *ListingService) Search(ctx context.Context, input dto.SearchListingsInput) (*dto.SearchListingsResponse, error) {
	if s.search == nil {
		return nil, fmt.Errorf("search client is not configured")
	}

	page := input.Page
	if page <= 0 {
		page = 1
	}
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	params := search.ListingSearchParams{
		Query:      input.Query,
		From:       (page - 1) * limit,
		Size:       limit,
		CategoryID: input.CategoryID,
		MinPrice:   input.MinPrice,
		MaxPrice:   input.MaxPrice,
	}

	result, err := s.search.SearchListings(ctx, params)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ListingResponse, 0, len(result.IDs))
	for _, id := range result.IDs {
		listing, err := s.listingRepo.GetByID(ctx, id)
		if err != nil {
			slog.Warn("failed to load listing from DB after search hit", "id", id, "error", err)
			continue
		}

		media, _ := s.mediaRepo.GetByListingID(ctx, id)
		chars, _ := s.charsRepo.GetByListingID(ctx, id)

		if resp := mapper.ToListingResponse(listing, chars, media); resp != nil {
			items = append(items, *resp)
		}
	}

	return &dto.SearchListingsResponse{
		Total: result.Total,
		Page:  page,
		Limit: limit,
		Items: items,
	}, nil
}

func (s *ListingService) enqueueIndexing(ctx context.Context, listing *domain.Listing, chars *domain.ListingCharacteristics) {
	if s.search == nil || listing == nil {
		return
	}

	doc := search.ListingDocument{
		ID:          listing.ID,
		UserID:      listing.UserID,
		CategoryID:  listing.CategoryID,
		Title:       listing.Title,
		Description: listing.Description,
		Price:       listing.Price,
		Currency:    string(listing.Currency),
		Status:      string(listing.Status),
		IsSold:      listing.IsSold,
		CreatedAt:   listing.CreatedAt.Format(time.RFC3339),
	}

	if chars != nil && len(chars.Tags) > 0 {
		doc.Tags = chars.Tags
	}

	go func() {
		if err := s.search.IndexListing(context.Background(), doc); err != nil {
			slog.Warn("failed to index listing", "id", listing.ID, "error", err)
		}
	}()
}

func (s *ListingService) enqueueDeletion(ctx context.Context, id uuid.UUID) {
	if s.search == nil {
		return
	}

	go func() {
		if err := s.search.DeleteListing(context.Background(), id); err != nil {
			slog.Warn("failed to delete listing from index", "id", id, "error", err)
		}
	}()
}
