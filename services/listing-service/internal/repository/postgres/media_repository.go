package pgrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"time"
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

func (r *MediaRepository) BulkCreate(ctx context.Context, mediaList []domain.ListingMedia) error {
	if len(mediaList) == 0 {
		return nil
	}

	query := `
		INSERT INTO listing_media (id, listing_id, file_url, file_type, mime_type, "order", created_at)
		VALUES (:id, :listing_id, :file_url, :file_type, :mime_type, :order, :created_at)
	`

	_, err := r.getQueryer(ctx).NamedExecContext(ctx, query, mediaList)
	if err != nil {
		return fmt.Errorf("failed to bulk create media: %w", err)
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

func (r *MediaRepository) GetMediaByListingIDs(ctx context.Context, listingIDs []uuid.UUID) (map[uuid.UUID][]domain.ListingMedia, error) {
	if len(listingIDs) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`
       SELECT id, listing_id, file_url, file_type, mime_type, "order", created_at
       FROM listing_media
       WHERE listing_id IN (?)
       ORDER BY listing_id, "order" ASC
    `, listingIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %w", err)
	}

	query = r.db.Rebind(query)

	var rows []domain.ListingMedia
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch listing media batch: %w", err)
	}

	result := make(map[uuid.UUID][]domain.ListingMedia)
	for _, m := range rows {
		result[m.ListingID] = append(result[m.ListingID], m)
	}

	return result, nil
}

func (r *MediaRepository) UpdateOrder(ctx context.Context, mediaID uuid.UUID, order int) error {
	query := `UPDATE listing_media SET order = $1 WHERE id = $2`
	_, err := r.getQueryer(ctx).ExecContext(ctx, query, order, mediaID)
	return err
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

func (r *MediaRepository) DeleteExpiredTemp(ctx context.Context, olderThan time.Duration) ([]TempMedia, error) {
	query := `
		SELECT * FROM media_uploads
		WHERE created_at < $1
		LIMIT 100
	`
	threshold := time.Now().Add(-olderThan)

	var batch []TempMedia
	if err := r.db.SelectContext(ctx, &batch, query, threshold); err != nil {
		return nil, err
	}

	if len(batch) == 0 {
		return nil, nil
	}

	ids := make([]uuid.UUID, len(batch))
	for i, m := range batch {
		ids[i] = m.ID
	}

	if err := r.DeleteTemp(ctx, ids); err != nil {
		return nil, err
	}

	return batch, nil
}

type TempMedia struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	S3Key     string    `db:"s3_key"`
	Filename  string    `db:"filename"`
	MimeType  string    `db:"mime_type"`
	Size      int64     `db:"size"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *MediaRepository) CreateTemp(ctx context.Context, media *TempMedia) error {
	query := `
        INSERT INTO media_uploads (id, user_id, s3_key, filename, mime_type, size, created_at)
        VALUES (:id, :user_id, :s3_key, :filename, :mime_type, :size, :created_at)
    `
	_, err := r.getQueryer(ctx).NamedExecContext(ctx, query, media)
	return err
}

func (r *MediaRepository) GetTempByIDs(ctx context.Context, ids []uuid.UUID) ([]TempMedia, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In("SELECT * FROM media_uploads WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var list []TempMedia
	err = r.db.SelectContext(ctx, &list, query, args...)
	return list, err
}

func (r *MediaRepository) DeleteTemp(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	query, args, err := sqlx.In("DELETE FROM media_uploads WHERE id IN (?)", ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)

	_, err = r.getQueryer(ctx).ExecContext(ctx, query, args...)
	return err
}
