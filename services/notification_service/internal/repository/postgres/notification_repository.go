package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/squ1ky/avigo-c2c-marketplace/services/notification-service/internal/domain"
)

var (
	ErrNotificationNotFound  = errors.New("notification not found")
	ErrInvalidNotificationID = errors.New("invalid notification ID")
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	if notification == nil {
		return errors.New("notification cannot be nil")
	}

	result := r.db.WithContext(ctx).Create(notification)
	if result.Error != nil {
		return fmt.Errorf("failed to create notification: %w", result.Error)
	}

	return nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	if id == uuid.Nil {
		return nil, ErrInvalidNotificationID
	}

	var notification domain.Notification
	result := r.db.WithContext(ctx).First(&notification, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotificationNotFound
		}
		return nil, fmt.Errorf("failed to get notification: %w", result.Error)
	}

	return &notification, nil
}

func (r *NotificationRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
	limit, offset int,
) ([]domain.Notification, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID cannot be empty")
	}

	var notifications []domain.Notification
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", result.Error)
	}

	return notifications, nil
}

func (r *NotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	if notification == nil {
		return errors.New("notification cannot be nil")
	}

	if notification.Id == uuid.Nil {
		return ErrInvalidNotificationID
	}

	result := r.db.WithContext(ctx).Save(notification)
	if result.Error != nil {
		return fmt.Errorf("failed to update notification: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotificationNotFound
	}

	return nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidNotificationID
	}

	now := time.Now()
	result := r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("id = ? AND read_at IS NULL", id).
		Update("read_at", now)

	if result.Error != nil {
		return fmt.Errorf("failed to mark notification as read: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotificationNotFound
	}

	return nil
}

func (r *NotificationRepository) CountUnreadByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	if userID == uuid.Nil {
		return 0, errors.New("user ID cannot be empty")
	}

	var count int64
	result := r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", result.Error)
	}

	return count, nil
}
