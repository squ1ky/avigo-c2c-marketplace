package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status string
type Role string

const (
	StatusPending Status = "pending"
	StatusActive  Status = "active"

	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username string    `gorm:"type:varchar(64);uniqueIndex;not null"`
	Status   Status    `gorm:"type:varchar(16);default:'pending'"`
	Role     Role      `gorm:"type:varchar(16);default:'user'"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Profile  *UserProfile  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
	Security *UserSecurity `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

type UserProfile struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id"`
	DisplayName string    `gorm:"type:varchar(64);not null" json:"display_name"`
	Phone       string    `gorm:"type:varchar(32)" json:"phone"`
	About       string    `gorm:"type:text" json:"about"`
	Country     string    `gorm:"type:varchar(64)" json:"country"`
	City        string    `gorm:"type:varchar(64)" json:"city"`
	AvatarURL   string    `gorm:"type:varchar(512)" json:"avatar_url"`

	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserSecurity struct {
	UserID                uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	Email                 string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash          string     `gorm:"type:varchar(255);not null" json:"-"`
	EmailVerified         bool       `gorm:"default:false" json:"email_verified"`
	RefreshToken          string     `gorm:"type:text" json:"-"`
	ConfirmationCode      string     `gorm:"type:text" json:"-"`
	ConfirmationExpiresAt *time.Time `json:"-"`
	LastLoginAt           *time.Time `json:"last_login_at,omitempty"`
}

func (*User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (*UserProfile) TableName() string {
	return "user_profile"
}

func (*UserSecurity) TableName() string {
	return "user_security"
}
