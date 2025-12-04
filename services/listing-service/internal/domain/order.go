package domain

import (
	"github.com/google/uuid"
	"time"
)

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "created"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID           uuid.UUID   `db:"id"`
	ListingID    uuid.UUID   `db:"listing_id"`
	BuyerID      uuid.UUID   `db:"buyer_id"`
	SellerID     uuid.UUID   `db:"seller_id"`
	Status       OrderStatus `db:"status"`
	CancelReason string      `db:"cancel_reason"`
	CreatedAt    time.Time   `db:"created_at"`
	UpdatedAt    time.Time   `db:"updated_at"`
}
