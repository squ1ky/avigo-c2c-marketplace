package dto

import (
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"time"
)

type OrderResponse struct {
	ID        uuid.UUID          `json:"id"`
	ListingID uuid.UUID          `json:"listing_id"`
	Status    domain.OrderStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`

	CounterpartyID   uuid.UUID `json:"counterparty_id"`
	CounterpartyName string    `json:"counterparty_name"`
	CounterpartyImg  string    `json:"counterparty_img"`
}
