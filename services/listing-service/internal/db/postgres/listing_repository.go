package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
)

type ListingRepository struct {
	db *sql.DB
}

func NewListingRepository(db *sql.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) Create(ctx context.Context, listing *domain.Listing) error {
	listing.ID = uuid.New().String()
	listing.CreatedAt = time.Now()
	listing.UpdatedAt = time.Now()

	metadataJSON, err := json.Marshal(listing.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	locationJSON, err := json.Marshal(listing.Location)
	if err != nil {
		return fmt.Errorf("failed to marshal location: %w", err)
	}

	query := `
		INSERT INTO listings (
			id, user_id, title, description, price, currency, category_id,
			status, photo_ids, video_ids, tags, views_count,
			location, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err = r.db.ExecContext(ctx, query,
		listing.ID,
		listing.UserID,
		listing.Title,
		listing.Description,
		listing.Price,
		listing.Currency,
		listing.CategoryID,
		listing.Status,
		pq.Array(listing.ID),
		pq.Array(listing.VideoIDs),
		pq.Array(listing.Tags),
		listing.ViewsCount,
		locationJSON,
		metadataJSON,
		listing.CreatedAt,
		listing.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert listing: %w", err)
	}

	return nil
}

func (r *ListingRepository) GetByID(ctx context.Context, id string) (*domain.Listing, error) {
	query := `
		SELECT
			id, user_id, title, description, price, currency, category_id,
			status, photo_ids, video_ids, tags, views_count,
			location, metadata, created_at, updated_at, published_at, expires_at,
		FROM listlings
		WHERE id = $1 AND status != $2
	`

	listing := &domain.Listing{}
	var locationJSON, metadataJSON []byte
	var publishedAt, expiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id, domain.StatusDeleted).Scan(
		&listing.ID,
		&listing.UserID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.Currency,
		&listing.CategoryID,
		&listing.Status,
		pq.Array(&listing.PhotoIDs),
		pq.Array(&listing.VideoIDs),
		pq.Array(&listing.Tags),
		&listing.ViewsCount,
		&locationJSON,
		&metadataJSON,
		&listing.CreatedAt,
		&listing.UpdatedAt,
		&publishedAt,
		&expiresAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("listing not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	if err := json.Unmarshal(locationJSON, &listing.Location); err != nil {
		return nil, fmt.Errorf("failed to unmarshal location: %w", err)
	}

	if err := json.Unmarshal(metadataJSON, &listing.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmmarshal metadata: %w", err)
	}
	if publishedAt.Valid {
		listing.PublishedAt = &publishedAt.Time
	}
	if expiresAt.Valid {
		listing.ExpiresAt = &expiresAt.Time
	}

	return listing, nil
}

func (r *ListingRepository) Update(ctx context.Context, listing *domain.Listing) error {
	listing.UpdatedAt = time.Now()

	metadataJSON, err := json.Marshal(listing.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	locationJSON, err := json.Marshal(listing.Location)
	if err != nil {
		return fmt.Errorf("failed to marshal location: %w", err)
	}

	query := `
        UPDATE listings SET
            title = $2, description = $3, price = $4, currency = $5,
            category_id = $6, photo_ids = $7, video_ids = $8, tags = $9,
            location = $10, metadata = $11, updated_at = $12
        WHERE id = $1
    `

	result, err := r.db.ExecContext(ctx, query,
		listing.ID,
		listing.Title,
		listing.Description,
		listing.Price,
		listing.Currency,
		listing.CategoryID,
		pq.Array(listing.PhotoIDs),
		pq.Array(listing.VideoIDs),
		pq.Array(listing.Tags),
		locationJSON,
		metadataJSON,
		listing.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update listing: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("listing not found")
	}

	return nil
}

func (r *ListingRepository) Delete(ctx context.Context, id string) error {
	query := `
        UPDATE listings SET
            status = $2, updated_at = $3
        WHERE id = $1
    `

	result, err := r.db.ExecContext(ctx, query, id, domain.StatusDeleted, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("listing not found")
	}

	return nil
}

func (r *ListingRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Listing, error) {
	query := `
        SELECT 
            id, user_id, title, description, price, currency, category_id,
            status, photo_ids, video_ids, tags, views_count,
            location, metadata, created_at, updated_at, published_at, expires_at
        FROM listings
        WHERE user_id = $1 AND status != $2
        ORDER BY created_at DESC
        LIMIT $3 OFFSET $4
    `

	rows, err := r.db.QueryContext(ctx, query, userID, domain.StatusDeleted, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query listings: %w", err)
	}
	defer rows.Close()

	return r.scanListings(rows)
}

func (r *ListingRepository) scanListings(rows *sql.Rows) ([]*domain.Listing, error) {
	var listings []*domain.Listing

	for rows.Next() {
		listing := &domain.Listing{}
		var locationJSON, metadataJSON []byte
		var publishedAt, expiresAt sql.NullTime

		err := rows.Scan(
			&listing.ID, &listing.UserID, &listing.Title, &listing.Description,
			&listing.Price, &listing.Currency, &listing.CategoryID, &listing.Status,
			pq.Array(&listing.PhotoIDs), pq.Array(&listing.VideoIDs), pq.Array(&listing.Tags),
			&listing.ViewsCount, &locationJSON, &metadataJSON,
			&listing.CreatedAt, &listing.UpdatedAt, &publishedAt, &expiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}

		if err := json.Unmarshal(locationJSON, &listing.Location); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(metadataJSON, &listing.Metadata); err != nil {
			return nil, err
		}

		if publishedAt.Valid {
			listing.PublishedAt = &publishedAt.Time
		}
		if expiresAt.Valid {
			listing.ExpiresAt = &expiresAt.Time
		}

		listings = append(listings, listing)
	}

	return listings, nil
}
