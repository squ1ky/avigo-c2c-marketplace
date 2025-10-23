package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/domain"
	model "github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/dto"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	GetByUserID(ctx context.Context, id uuid.UUID, limit, offset int) ([]domain.Notification, error)
	Update(ctx context.Context, notification *domain.Notification) error
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	CountUnreadByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

type EmailSender interface {
	Send(ctx context.Context, message *domain.EmailMessage) error
}

type TemplateData map[string]interface{}

type TemplatRender interface {
	Render(ctx context.Context, templateName string, data TemplateData) (string, error)
}

type EmailNotificationUseCase struct {
	notificationRepo NotificationRepository
	emailSender      EmailSender
	templateRender   TemplatRender
}

func NewEmailNotificationUseCase(
	repo NotificationRepository,
	emailSender EmailSender,
	render TemplatRender,
) *EmailNotificationUseCase {
	return &EmailNotificationUseCase{
		notificationRepo: repo,
		emailSender:      emailSender,
		templateRender:   render,
	}
}

func (uc *EmailNotificationUseCase) SendRegistrationEmail(
	ctx context.Context,
	req model.SendEmailNotificationRequest,
) (*model.SendEmailNotificationResponse, error) {

	if err := req.Validate(); err != nil {
		return &model.SendEmailNotificationResponse{
			Status: "failed",
			Error:  fmt.Errorf("validation failed: %w", err),
		}, err
	}

	notification := &domain.Notification{
		UserID:  req.UserID,
		Type:    domain.TypeEmailVerification,
		Title:   "Подтверждение регистрации",
		Message: fmt.Sprintf("Ваш код подтверждения: %s", req.ConfirmationCode),
		Metadata: domain.JSONB{
			"confirmation_code": req.ConfirmationCode,
			"email":             req.Email,
			"display_name":      req.DisplayName,
			"expires_at":        req.ExpiresAt.Format(time.RFC3339),
		},
		Status: domain.StatusPending,
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return &model.SendEmailNotificationResponse{
			Status: "failed",
			Error:  fmt.Errorf("failed to create notification: %w", err),
		}, err
	}

	expiresIn := time.Until(req.ExpiresAt).Minutes()
	templateData := TemplateData{
		"DisplayName":      req.DisplayName,
		"ConfirmationCode": req.ConfirmationCode,
		"ExpiresIn":        fmt.Sprintf("%.0f минут", expiresIn),
	}

	htmlBody, err := uc.templateRender.Render(ctx, "registration", templateData)
	if err != nil {
		_ = uc.notificationRepo.Update(ctx, notification)
		return &model.SendEmailNotificationResponse{
			NotificationID: notification.Id,
			Status:         "failde",
			Error:          fmt.Errorf("failed to render template: %w", err),
		}, err
	}

	emailMsg := &domain.EmailMessage{
		To:       req.Email,
		Subject:  "Подтверждение регистрации в АвиGo Маркетплейс",
		HTMLBody: htmlBody,
		TextBody: "",
	}

	if err := uc.emailSender.Send(ctx, emailMsg); err != nil {
		_ = uc.notificationRepo.Update(ctx, notification)
		return &model.SendEmailNotificationResponse{
			NotificationID: notification.Id,
			Status:         "failed",
			Error:          fmt.Errorf("failed to send email: %w", err),
		}, err
	}

	return &model.SendEmailNotificationResponse{
		NotificationID: notification.Id,
		Status:         "sent",
		SentAt:         time.Now(),
		Error:          nil,
	}, nil
}
