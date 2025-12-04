package domain

import (
	"github.com/google/uuid"
	"time"
)

type Review struct {
	ID             uuid.UUID `db:"id"`
	OrderID        uuid.UUID `db:"order_id"`
	ReviewerID     uuid.UUID `db:"reviewer_id"`      // Buyer
	ReviewedUserID uuid.UUID `db:"reviewed_user_id"` // Seller
	Rating         int       `db:"rating"`
	Text           string    `db:"text"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
