package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateReviewInput struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
	Rating  int       `json:"rating" binding:"required,min=1,max=5"` // 1-5
	Text    string    `json:"text"`

	// We will receive this from JWT
	ReviewerID uuid.UUID `json:"-"`
}

type ReviewResponse struct {
	ID         uuid.UUID `json:"id"`
	ReviewerID uuid.UUID `json:"reviewer_id"`
	Rating     int       `json:"rating"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}
