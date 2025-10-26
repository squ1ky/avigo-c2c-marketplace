package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/domain"
)

type SendEmailNotificationRequest struct {
	UserID           uuid.UUID
	Email            string
	DisplayName      string
	ConfirmationCode string
	ExpiresAt        time.Time
}

type SendEmailNotificationResponse struct {
	NotificationID uuid.UUID
	Status         string
	SentAt         time.Time
	Error          error
}

func (r *SendEmailNotificationResponse) IsSuccess() bool {
	return r.Error == nil
}

func (r *SendEmailNotificationRequest) Validate() error {
	if r.UserID == uuid.Nil {
		return domain.ErrInvalidUUID
	}
	if r.Email == "" {
		return domain.ErrMissingEmail
	}
	if r.ConfirmationCode == "" {
		return domain.ErrMissingConfirmationCode
	}
	if r.ExpiresAt.Before(time.Now()) {
		return domain.ErrExpiredCode
	}
	return nil
}
