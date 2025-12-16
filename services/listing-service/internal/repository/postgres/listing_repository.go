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

type ListingRepository struct {
	db *sqlx.DB
}

func NewListingRepository(db *sqlx.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) getQueryer(ctx context.Context) SQLQueryer {
	if tx := injectTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *ListingRepository) Create(ctx context.Context, listing *domain.Listing) error {
	query := `
		INSERT INTO listings (id, user_id, category_id, title, description, price, currency, status, views_count, is_sold, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	listing.ViewsCount = 0
	listing.IsSold = false
	listing.Status = domain.ListingStatusActive

	_, err := r.getQueryer(ctx).ExecContext(ctx, query,
		listing.ID,
		listing.UserID,
		listing.CategoryID,
		listing.Title,
		listing.Description,
		listing.Price,
		listing.Currency,
		listing.Status,
		listing.ViewsCount,
		listing.IsSold,
		listing.CreatedAt,
		listing.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create listing: %w", err)
	}

	return nil
}

func (r *ListingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Listing, error) {
	query := `
		SELECT id, user_id, category_id, title, description, price, currency, status, views_count, is_sold, created_at, updated_at
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

func (r *ListingRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*domain.Listing, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`
       SELECT id, user_id, category_id, title, description, price, currency, status, views_count, is_sold, created_at, updated_at
       FROM listings
       WHERE id IN (?)
    `, ids)

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	query = r.db.Rebind(query)

	var listings []domain.Listing
	if err := r.db.SelectContext(ctx, &listings, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch listings batch: %w", err)
	}

	result := make(map[uuid.UUID]*domain.Listing)
	for i := range listings {
		result[listings[i].ID] = &listings[i]
	}

	return result, nil
}

func (r *ListingRepository) Update(ctx context.Context, listing *domain.Listing) error {
	query := `
		UPDATE listings
		SET title = :title,
		    description = :description,
		    price = :price,
		    currency = :currency,
		    status = :status,
		    category_id = :category_id,
		    updated_at = :updated_at
		WHERE id = :id
	`

	listing.UpdatedAt = time.Now()

	result, err := r.getQueryer(ctx).NamedExecContext(ctx, query, listing)
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

	result, err := r.getQueryer(ctx).ExecContext(ctx, query, id)
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

func (r *ListingRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.ListingWithMedia, error) {
	query := `
		SELECT l.id, l.user_id, l.category_id, l.title, l.description,
		       l.price, l.currency, l.status, l.views_count, l.is_sold, l.created_at, l.updated_at,
			   COALESCE(m.file_url, '') as main_image_url
		FROM listings l
		LEFT JOIN listing_media m on l.id = m.listing_id AND m."order" = 0
		WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var rows []struct {
		domain.Listing
		MainImageURL string `db:"main_image_url"`
	}

	if err := r.db.SelectContext(ctx, &rows, query, userID, limit, offset); err != nil {
		return nil, fmt.Errorf("failed to get listings by user: %w", err)
	}

	results := make([]domain.ListingWithMedia, len(rows))
	for i, row := range rows {
		results[i] = domain.ListingWithMedia{
			Listing:   row.Listing,
			MainImage: row.MainImageURL,
		}
	}

	return results, nil
}

func (r *ListingRepository) IncrementViews(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE listings SET views_count = views_count + 1 WHERE id = $1`

	_, err := r.getQueryer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	return nil
}

func (r *ListingRepository) GetByCategoryID(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]domain.Listing, error) {
	query := `
		SELECT id, user_id, category_id, title, description, price, currency, status, views_count, is_sold, created_at, updated_at
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

func (r *ListingRepository) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `
		SELECT id, name, slug, parent_id, level, "order"
		FROM categories
		ORDER BY level ASC, "order" ASC
	`

	var categories []*domain.Category
	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

func (r *ListingRepository) MarkAsSold(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE listings
		SET is_sold = true, status = 'inactive', updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.getQueryer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark listing as sold: %w", err)
	}

	return nil
}

func (r *ListingRepository) MarkAsActive(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE listings
		SET is_sold = false, status = 'active', updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.getQueryer(ctx).ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark listing as active: %w", err)
	}

	if _, err := result.RowsAffected(); err != nil {
		return nil
	}

	return nil
}
