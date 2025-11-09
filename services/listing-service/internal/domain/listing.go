package domain

import (
	"github.com/google/uuid"
	"time"
)

type Listing struct {
	ID          uuid.UUID     `db:"id"`
	UserID      uuid.UUID     `db:"user_id"`
	CategoryID  uuid.UUID     `db:"category_id"`
	Title       string        `db:"title"`
	Description string        `db:"description"`
	Price       float64       `db:"price"`
	Currency    Currency      `db:"currency"`
	Status      ListingStatus `db:"status"`
	ViewsCount  int64         `db:"views_count"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}

type ListingStatus string

const (
	ListingStatusDraft    ListingStatus = "draft"
	ListingStatusActive   ListingStatus = "active"
	ListingStatusInactive ListingStatus = "inactive"
)

type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

type Category struct {
	ID        uuid.UUID  `db:"id"`
	Name      string     `db:"name"`
	Slug      string     `db:"slug"`      // for URL
	ParentID  *uuid.UUID `db:"parent_id"` // NULL for root category
	Level     int        `db:"level"`     // 0 - electronics, 1 - notebooks, 2 - for work
	Order     int        `db:"order"`
	CreatedAt time.Time  `db:"created_at"`
}

type ListingCharacteristics struct {
	ListingID       string                 `bson:"listing_id"`
	Characteristics map[string]interface{} `bson:"characteristics"`
	Tags            []string               `bson:"tags"`
	UpdatedAt       time.Time              `bson:"updated_at"`
}

type Favorite struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	ListingID uuid.UUID `db:"listing_id"`
	CreatedAt time.Time `db:"created_at"`
}
