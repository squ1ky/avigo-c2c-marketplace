package kafka

import "time"

type EventType string

const (
	EventTypeEmailVerification EventType = "user.email.verification.requested"
	// реализуем в будущем, по сути)
	EventTypeChatMessage    EventType = "user.notification.chat.message"
	EventTypeListingUpdate  EventType = "user.notification.listing.update"
	EventTypeReviewReceived EventType = "user.notification.review.received"
)

type UserEvent struct {
	EventType EventType `json:"event_type"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type EmailVerificationEvent struct {
	UserEvent
	Email            string    `json:"email"`
	DisplayName      string    `json:"display_name"`
	ConfirmationCode string    `json:"confirmation_code"`
	ExpiresAt        time.Time `json:"expires_at"`
}
