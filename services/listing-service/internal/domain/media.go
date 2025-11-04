package domain

import (
	"github.com/google/uuid"
	"time"
)

type MediaType string

const (
	MediaTypePhoto MediaType = "photo"
	MediaTypeVideo MediaType = "video"
)

type Media struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ListingID uuid.UUID `json:"listing_id" db:"listing_id"`
	Type      MediaType `json:"type" db:"type"`
	ObjectKey string    `json:"object_key" db:"object_key"` // MinIO path
	Order     int       `json:"order" db:"order"`
	Size      int64     `json:"size" db:"size"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
