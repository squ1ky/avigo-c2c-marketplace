package domain

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id     uuid.UUID
	UserID uuid.UUID
	Type   NotificationType

	Title   string
	Message string

	Metadata JSONB
	Status   NotificationStatus

	ReadAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (*Notification) TableName() string {
	return "notifications"
}

type NotificationType string

const (
	TypeEmailVerification NotificationType = "email_verification"
	TypeOrderCreated      NotificationType = "order_created"
	TypeOrderStatusChange NotificationType = "order_status_change"
	TypeNewMessage        NotificationType = "new_message"
	TypeReviewCreated     NotificationType = "review_created"
	TypeNewReview         NotificationType = "new_review"
)

type DeliveryChannel string

const (
	ChannelEmail DeliveryChannel = "email"
	ChannelPush  DeliveryChannel = "push"
)

type NotificationStatus string

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
)
