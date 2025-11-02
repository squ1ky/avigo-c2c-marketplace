package domain

import (
	"time"
)

type Listing struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Currency    string     `json:"currency"`
	CategoryID  string     `json:"category_id"`
	Status      string     `json:"status"` // draft, pending_moderation, published, archived, deleted
	PhotoIDs    []string   `json:"photo_ids"`
	VideoIDs    []string   `json:"video_ids,omitempty"`
	Tags        []string   `json:"tags"`
	ViewsCount  int64      `json:"views_count"`
	Location    Location   `json:"location"`
	Metadata    Metadata   `json:"metadata"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type Location struct {
	City      string  `json:"city"`
	Region    string  `json:"region,omitempty"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type Metadata struct {
	ResponsesCount int64                  `json:"responses_count"`
	FavoritesCount int64                  `json:"favorites_count"`
	CustomFields   map[string]interface{} `json:"custom_fields,omitempty"`
}

type ListingStatus string

const (
	StatusDraft             ListingStatus = "draft"
	StatusPendingModeration ListingStatus = "pending_moderation"
	StatusPublished         ListingStatus = "published"
	StatusArchived          ListingStatus = "archived"
	StatusDeleted           ListingStatus = "deleted"
	StatusRejected          ListingStatus = "rejected"
)

type Photo struct {
	ID          string    `json:"id"`
	ListingID   string    `json:"listing_id"`
	ObjectID    string    `json:"object_id"` // ID в S3/MinIO
	URL         string    `json:"url"`       // Presigned URL
	Order       int       `json:"order"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}

type Video struct {
	ID          string    `json:"id"`
	ListingID   string    `json:"listing_id"`
	ObjectID    string    `json:"object_id"`
	URL         string    `json:"url"`
	ThumbnailID string    `json:"thumbnail_id,omitempty"`
	Duration    int       `json:"duration"` // секунды
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}
