package domain

import (
	"github.com/google/uuid"
	"time"
)

type Listing struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Amount      float64    `json:"amount" db:"amount"`
	Currency    string     `json:"currency" db:"currency"`
	CategoryID  uuid.UUID  `json:"category_id" db:"category_id"`
	Country     string     `json:"country" db:"country"`
	City        string     `json:"city" db:"city"`
	ViewsCount  int64      `json:"views_count" db:"views_count"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}

type ListingMetadata struct {
	ListingID      uuid.UUID              `bson:"listing_id" json:"listing_id"`
	ResponsesCount int64                  `bson:"responses_count" json:"responses_count"`
	FavoritesCount int64                  `bson:"favorites_count" json:"favorites_count"`
	SearchKeywords []string               `bson:"search_keywords" json:"search_keywords"`
	Attributes     map[string]interface{} `bson:"attributes" json:"attributes"`
	CreatedAt      time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time              `bson:"updated_at" json:"updated_at"`
}
