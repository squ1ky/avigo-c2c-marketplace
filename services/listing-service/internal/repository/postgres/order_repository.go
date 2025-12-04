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

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) getQueryer(ctx context.Context) SQLQueryer {
	if tx := injectTx(ctx); tx != nil {
		return tx
	}
	return r.db
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `
        INSERT INTO orders (id, listing_id, buyer_id, seller_id, status, cancel_reason, created_at, updated_at)
        VALUES (:id, :listing_id, :buyer_id, :seller_id, :status, :cancel_reason, :created_at, :updated_at)
    `

	order.ID = uuid.New()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	_, err := r.getQueryer(ctx).NamedExecContext(ctx, query, order)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `
        SELECT id, listing_id, buyer_id, seller_id, status, cancel_reason, created_at, updated_at
        FROM orders
        WHERE id = $1
    `

	var order domain.Order
	if err := r.db.GetContext(ctx, &order, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus, cancelReason string) error {
	query := `
        UPDATE orders 
        SET status = $1, cancel_reason = $2, updated_at = NOW() 
        WHERE id = $3
    `

	_, err := r.getQueryer(ctx).ExecContext(ctx, query, status, cancelReason, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (r *OrderRepository) GetPurchasesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	query := `
        SELECT id, listing_id, buyer_id, seller_id, status, cancel_reason, created_at, updated_at
        FROM orders
        WHERE buyer_id = $1
        ORDER BY created_at DESC
    `

	var orders []domain.Order
	if err := r.db.SelectContext(ctx, &orders, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get purchases: %w", err)
	}

	return orders, nil
}

func (r *OrderRepository) GetSalesByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Order, error) {
	query := `
        SELECT id, listing_id, buyer_id, seller_id, status, cancel_reason, created_at, updated_at
        FROM orders
        WHERE seller_id = $1
        ORDER BY created_at DESC
    `

	var orders []domain.Order
	if err := r.db.SelectContext(ctx, &orders, query, userID); err != nil {
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}

	return orders, nil
}
