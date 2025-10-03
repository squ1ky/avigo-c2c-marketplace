package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusActive  Status = "active"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Username string    `gorm:"type:varchar(64);uniqueIndex;not null"`
	Status   Status    `gorm:"type:varchar(16);default:'pending'"`
	Role     string    `gorm:"type:varchar(16);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Profile  *UserProfile  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	Security *UserSecurity `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" `
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
