package dto

import (
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"time"
)

type OrderResponse struct {
	ID           uuid.UUID            `json:"id"`
	Status       domain.OrderStatus   `json:"status"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	Listing      ListingShortResponse `json:"listing"`
	Counterparty UserResponse         `json:"counterparty"`
}
