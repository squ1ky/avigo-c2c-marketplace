package dto

import (
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"io"
	"time"
)

type CreateListingInput struct {
	UserID          uuid.UUID
	CategoryID      uuid.UUID
	Title           string
	Description     string
	Price           float64
	Currency        domain.Currency
	Files           []FileInput
	Characteristics map[string]interface{}
	Tags            []string
}

type ListingResponse struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.UUID              `json:"user_id"`
	CategoryID      uuid.UUID              `json:"category_id"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Price           float64                `json:"price"`
	Currency        domain.Currency        `json:"currency"`
	Status          domain.ListingStatus   `json:"status"`
	ViewsCount      int64                  `json:"views_count"`
	IsSold          bool                   `json:"is_sold"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Characteristics map[string]interface{} `json:"characteristics"`
	Tags            []string               `json:"tags"`
	Media           []domain.ListingMedia  `json:"media"`
}

type FileInput struct {
	Name        string
	ContentType string
	Size        int64
	Data        io.Reader
}
