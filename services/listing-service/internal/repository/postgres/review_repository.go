package pgrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
)

type ReviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) getQueryer(ctx context.Context) SQLQueryer {
	if tx := injectTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *ReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	query := `
        INSERT INTO reviews (id, order_id, reviewer_id, reviewed_user_id, rating, text, created_at, updated_at)
        VALUES (:id, :order_id, :reviewer_id, :reviewed_user_id, :rating, :text, :created_at, :updated_at)
    `

	review.ID = uuid.New()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	_, err := r.getQueryer(ctx).NamedExecContext(ctx, query, review)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

func (r *ReviewRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Review, error) {
	query := `
        SELECT id, order_id, reviewer_id, reviewed_user_id, rating, text, created_at, updated_at
        FROM reviews
        WHERE reviewed_user_id = $1
        ORDER BY created_at DESC
    `

	var reviews []domain.Review
	if err := r.db.SelectContext(ctx, &reviews, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}

func (r *ReviewRepository) ExistsByOrderID(ctx context.Context, orderID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM reviews WHERE order_id = $1)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, orderID).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check review existence: %w", err)
	}

	return exists, nil
}
