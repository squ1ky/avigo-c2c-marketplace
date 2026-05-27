package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/dto"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/grpc/client/user"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/mapper"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/mongo"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/postgres"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/repository/s3"
	"log/slog"
	"time"
)

type ListingService struct {
	listingRepo *pgrepo.ListingRepository
	mediaRepo   *pgrepo.MediaRepository
	charsRepo   *mngrepo.CharacteristicsRepository
	storage     *s3.MediaStorage
	txManager   pgrepo.TransactionManager
	s3Config    config.S3Config
	userClient  *user.Client
}

func NewListingService(
	listingRepo *pgrepo.ListingRepository,
	mediaRepo *pgrepo.MediaRepository,
	charsRepo *mngrepo.CharacteristicsRepository,
	storage *s3.MediaStorage,
	txManager pgrepo.TransactionManager,
	s3Config config.S3Config,
	userClient *user.Client,
) *ListingService {
	return &ListingService{
		listingRepo: listingRepo,
		mediaRepo:   mediaRepo,
		charsRepo:   charsRepo,
		storage:     storage,
		txManager:   txManager,
		s3Config:    s3Config,
		userClient:  userClient,
	}
}

func (s *ListingService) Create(ctx context.Context, input dto.CreateListingInput) (*dto.ListingResponse, error) {
	listingID := uuid.New()
	now := time.Now()
	input.Tags = normalizeTags(input.Tags, 20)

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

	return mapper.ToListingResponse(listing, chars, nil, s.s3Config.PublicURL), nil
}

func (s *ListingService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ListingWithUserResponse, error) {
	listing, err := s.listingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	media, _ := s.mediaRepo.GetByListingID(ctx, id)
	chars, _ := s.charsRepo.GetByListingID(ctx, id)

	userCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	userGrpc, err := s.userClient.GetUserByID(userCtx, listing.UserID.String())
	if err != nil {
		slog.Error("failed to fetch user, degrading response", "error", err)
		userGrpc = nil
	}

	go func() {
		_ = s.listingRepo.IncrementViews(context.Background(), id)
	}()

	return mapper.ToListingWithUserResponse(listing, chars, media, userGrpc, s.s3Config.PublicURL), nil
}

func (s *ListingService) GetCategoriesTree(ctx context.Context) ([]*dto.CategoryResponse, error) {
	categories, err := s.listingRepo.GetAllCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return s.buildCategoryTree(categories), nil
}

func (s *ListingService) buildCategoryTree(allCats []*domain.Category) []*dto.CategoryResponse {
	categoryMap := make(map[string]*dto.CategoryResponse)
	var roots []*dto.CategoryResponse

	for _, cat := range allCats {
		categoryMap[cat.ID.String()] = &dto.CategoryResponse{
			ID:       cat.ID.String(),
			Name:     cat.Name,
			Slug:     cat.Slug,
			Children: make([]*dto.CategoryResponse, 0),
		}
		if cat.ParentID != nil {
			pid := cat.ParentID.String()
			categoryMap[cat.ID.String()].ParentID = &pid
		}
	}

	for _, cat := range allCats {
		dtoCat := categoryMap[cat.ID.String()]

		if cat.ParentID == nil {
			roots = append(roots, dtoCat)
		} else {
			parentID := cat.ParentID.String()
			if parent, exists := categoryMap[parentID]; exists {
				parent.Children = append(parent.Children, dtoCat)
			}
		}
	}

	return roots
}

func (s *ListingService) GetUserListings(ctx context.Context, userID uuid.UUID, limit, offset int) ([]dto.ListingResponse, error) {
	listings, err := s.listingRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user listings: %w", err)
	}

	responses := make([]dto.ListingResponse, 0, len(listings))

	for _, listing := range listings {
		var mediaList []domain.ListingMedia

		if listing.MainImage != "" {
			mediaList = append(mediaList, domain.ListingMedia{
				ListingID: listing.ID,
				FileURL:   listing.MainImage,
				FileType:  domain.MediaTypeImage,
				Order:     0,
			})
		}

		resp := mapper.ToListingResponse(&listing.Listing, nil, mediaList, s.s3Config.PublicURL)
		responses = append(responses, *resp)
	}

	return responses, nil
}

func (s *ListingService) GetCatalogListings(ctx context.Context, limit, offset int) ([]dto.ListingResponse, error) {
	listings, err := s.listingRepo.GetAllActive(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch catalog listings: %w", err)
	}

	responses := make([]dto.ListingResponse, 0, len(listings))
	for _, listing := range listings {
		var mediaList []domain.ListingMedia

		if listing.MainImage != "" {
			mediaList = append(mediaList, domain.ListingMedia{
				ListingID: listing.ID,
				FileURL:   listing.MainImage,
				FileType:  domain.MediaTypeImage,
				Order:     0,
			})
		}

		resp := mapper.ToListingResponse(&listing.Listing, nil, mediaList, s.s3Config.PublicURL)
		responses = append(responses, *resp)
	}

	return responses, nil
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
	input.Tags = normalizeTags(input.Tags, 20)
	chars := mapper.ToDomainCharacteristicsFromUpdate(input, existing.UpdatedAt)

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

	if err := s.charsRepo.Upsert(ctx, chars); err != nil {
		return nil, fmt.Errorf("failed to update listing characteristics: %w", err)
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

	return mapper.ToListingResponse(existing, chars, finalMediaList, s.s3Config.PublicURL), nil
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
				slog.Warn("failed to cleanup s3 files", "error", err)
			}
		}()
	}

	return nil
}
