package pgrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
)

type MediaRepository struct {
	db *sqlx.DB
}

func NewMediaRepository(db *sqlx.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) getQueryer(ctx context.Context) SQLQueryer {
	if tx := injectTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *MediaRepository) Create(ctx context.Context, media *domain.ListingMedia) error {
	query := `
		INSERT INTO listing_media (id, listing_id, file_url, file_type, mime_type, "order", created_at)
		VALUES (:id, :listing_id, :file_url, :file_type, :mime_type, :order, :created_at)
	`

	media.ID = uuid.New()

	_, err := r.getQueryer(ctx).NamedExecContext(ctx, query, media)
	if err != nil {
		return fmt.Errorf("failed to create listing media: %w", err)
	}

	return nil
}

func (r *MediaRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ListingMedia, error) {
	query := `
		SELECT id, listing_id, file_url, file_type, mime_type, "order", created_at
    	FROM listing_media
    	WHERE id = $1
	`

	var media domain.ListingMedia
	if err := r.db.GetContext(ctx, &media, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("media not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get media: %w", err)
	}

	return &media, nil
}

func (r *MediaRepository) GetByListingID(ctx context.Context, listingID uuid.UUID) ([]domain.ListingMedia, error) {
	query := `
		SELECT id, listing_id, file_url, file_type, mime_type, "order", created_at
		FROM listing_media
		WHERE listing_id = $1
		ORDER BY "order" ASC
	`

	var media []domain.ListingMedia
	if err := r.db.SelectContext(ctx, &media, query, listingID); err != nil {
		return nil, fmt.Errorf("failed to get listing media: %w", err)
	}

	return media, nil
}

func (r *MediaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM listing_media WHERE id = $1`

	_, err := r.getQueryer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete listing media: %w", err)
	}

	return nil
}

func (r *MediaRepository) DeleteAllByListingID(ctx context.Context, listingID uuid.UUID) error {
	query := `DELETE FROM listing_media WHERE listing_id = $1`

	_, err := r.getQueryer(ctx).ExecContext(ctx, query, listingID)
	if err != nil {
		return fmt.Errorf("failed to delete listing media: %w", err)
	}

	return nil
}
