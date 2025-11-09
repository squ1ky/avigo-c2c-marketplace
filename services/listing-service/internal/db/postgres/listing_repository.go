package postgres

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

type ListingRepository struct {
	db *sqlx.DB
}

func NewListingRepository(db *sqlx.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) Create(ctx context.Context, listing *domain.Listing) error {
	query := `
		INSERT INTO listings (id, user_id, category_id, title, description, price, currency, status, views_count, created_at, updated_at)
		VALUES (:id, :user_id, :category_id, :title, :description, :price, :currency, :status, :views_count, :created_at, :updated_at)
	`

	listing.ID = uuid.New()
	listing.CreatedAt = time.Now()
	listing.UpdatedAt = time.Now()
	listing.ViewsCount = 0
	listing.Status = domain.ListingStatusActive

	_, err := r.db.ExecContext(ctx, query, listing)
	if err != nil {
		return fmt.Errorf("failed to create listing: %w", err)
	}

	return nil
}

func (r *ListingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Listing, error) {
	query := `
		SELECT id, user_id, category_id, title, description, price, currency, status, views_count, created_at, updated_at
		FROM listings
		WHERE id = $1
	`

	var listing domain.Listing
	err := r.db.GetContext(ctx, &listing, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("listing not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	return &listing, nil
}

func (r *ListingRepository) Update(ctx context.Context, listing *domain.Listing) error {
	query := `
		UPDATE listings
		SET title = :title,
		    description = :description,
		    price = :price,
		    currency = :currency,
		    status = :status,
		    category = :category_id,
		    updated_at = :updated_at
		WHERE id = :id
	`

	listing.UpdatedAt = time.Now()

	result, err := r.db.NamedExecContext(ctx, query, listing)
	if err != nil {
		return fmt.Errorf("failed to update listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing not found")
	}

	return nil
}

func (r *ListingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM listings WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("listing not found")
	}

	return nil
}

func (r *ListingRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Listing, error) {
	query := `
		SELECT id, user_id, category_id, title, description, price, currency, status, views_count, created_at, updated_at
		FROM listings
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var listings []domain.Listing
	if err := r.db.SelectContext(ctx, &listings, query, userID, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get listings by user: %w", err)
	}

	return listings, nil
}

func (r *ListingRepository) IncrementViews(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE listings SET views_count = views_count + 1 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	return nil
}

func (r *ListingRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]domain.Listing, error) {
	query := `
		SELECT id, user_id, category_id, title, description, price, currency, status, views_count, created_at, updated_at
		FROM listings
		WHERE category_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var listings []domain.Listing
	if err := r.db.SelectContext(ctx, &listings, query, categoryID, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get listings by category: %w", err)
	}

	return listings, nil
}
