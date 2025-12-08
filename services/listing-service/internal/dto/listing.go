package dto

import (
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/domain"
	"time"
)

type CreateListingInput struct {
	UserID          uuid.UUID              `json:"-"`
	CategoryID      uuid.UUID              `json:"category_id" binding:"required"`
	Title           string                 `json:"title" binding:"required,min=5,max=100"`
	Description     string                 `json:"description" binding:"required,min=10,max=2000"`
	Price           float64                `json:"price" binding:"required,min=0"`
	Currency        domain.Currency        `json:"currency" binding:"required,oneof=RUB USD EUR"`
	MediaIDs        []uuid.UUID            `json:"media_ids" binding:"omitempty,dive,uuid"`
	Characteristics map[string]interface{} `json:"characteristics"`
	Tags            []string               `json:"tags"`
}

type UpdateListingInput struct {
	UserID          uuid.UUID              `json:"-"`
	ID              uuid.UUID              `json:"-"`
	CategoryID      uuid.UUID              `json:"category_id" binding:"required"`
	Title           string                 `json:"title" binding:"required,min=5,max=100"`
	Description     string                 `json:"description" binding:"required,min=10,max=2000"`
	Price           float64                `json:"price" binding:"required,min=0"`
	Currency        domain.Currency        `json:"currency" binding:"required,oneof=RUB USD EUR"`
	MediaIDs        []uuid.UUID            `json:"media_ids" binding:"omitempty,dive,uuid"`
	Characteristics map[string]interface{} `json:"characteristics"`
	Tags            []string               `json:"tags"`
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
