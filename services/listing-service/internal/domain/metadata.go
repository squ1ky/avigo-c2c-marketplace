package domain

import "time"

type Category struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Slug        string    `bson:"slug" json:"slug"`
	ParentID    *string   `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	Icon        string    `bson:"icon,omitempty" json:"icon,omitempty"`
	Order       int       `bson:"order" json:"order"`
	IsActive    bool      `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

type Tag struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Slug      string    `bson:"slug" json:"slug"`
	Type      string    `bson:"type" json:"type"`
	IsActive  bool      `bson:"is_active" json:"is_active"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type ListingMetadata struct {
	ID              string                 `bson:"_id,omitempty" json:"id"`
	ListingID       string                 `bson:"listing_id" json:"listing_id"`
	SearchKeywords  []string               `bson:"search_keywords" json:"search_keywords"`
	Attributes      map[string]interface{} `bson:"attributes" json:"attributes"`
	ModerationNotes string                 `bson:"moderation_notes,omitempty" json:"moderation_notes,omitempty"`
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`
}
