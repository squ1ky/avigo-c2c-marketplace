package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	model "github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/dto"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/usecase"
)

type NotificationHandler struct {
	emailUseCase *usecase.EmailNotificationUseCase
}

func NewNotificationHandler(emailUseCase *usecase.EmailNotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		emailUseCase: emailUseCase,
	}
}

func (h *NotificationHandler) Handle(ctx context.Context, message kafka.Message) error {
	log.Printf("Received message: key=%s, value=%s", string(message.Key), string(message.Value))

	var baseEvent UserEvent
	if err := json.Unmarshal(message.Value, &baseEvent); err != nil {
		return fmt.Errorf("failed to unmarshal base event: %w", err)
	}

	switch baseEvent.EventType {
	case EventTypeEmailVerification:
		return h.handleEmailVerification(ctx, message.Value)
	default:
		log.Printf("Unknown event type: %s", baseEvent.EventType)
		return nil
	}
}

func (h *NotificationHandler) handleEmailVerification(ctx context.Context, data []byte) error {

	var event EmailVerificationEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to unmarshal email verification event: %w", err)
	}

	// if err := event.Validate(); err != nil {
	//     return fmt.Errorf("event validation failed: %w", err)
	// }

	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	req := model.SendEmailNotificationRequest{
		UserID:           userID,
		Email:            event.Email,
		DisplayName:      event.DisplayName,
		ConfirmationCode: event.ConfirmationCode,
		ExpiresAt:        event.ExpiresAt,
	}

	resp, err := h.emailUseCase.SendRegistrationEmail(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to send registration email: %w", err)
	}

	if !resp.IsSuccess() {
		return fmt.Errorf("email sending failed: %w", resp.Error)
	}

	log.Printf("Successfully sent registration email to %s (notification_id: %s)",
		req.Email, resp.NotificationID)

	return nil
}
