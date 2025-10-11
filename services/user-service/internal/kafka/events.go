package kafka

import "time"

const (
	EventTypeEmailVerification = "user.email.verification.requested"
)

type UserEvent struct {
	EventType string    `json:"event_type"`
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

func NewEmailVerificationEvent(userID, email, displayName, confirmationCode string, expiresAt time.Time) *EmailVerificationEvent {
	return &EmailVerificationEvent{
		UserEvent: UserEvent{
			EventType: EventTypeEmailVerification,
			UserID:    userID,
			Timestamp: time.Now().UTC(),
		},
		Email:            email,
		DisplayName:      displayName,
		ConfirmationCode: confirmationCode,
		ExpiresAt:        expiresAt,
	}
}
