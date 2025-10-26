package domain

import (
	"time"

	"github.com/google/uuid"
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

	Profile  *UserProfile  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Security *UserSecurity `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type UserProfile struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	DisplayName string    `gorm:"type:varchar(64);not null"`
	Phone       string    `gorm:"type:varchar(32)"`
	About       string    `gorm:"type:text"`
	Country     string    `gorm:"type:varchar(64)"`
	City        string    `gorm:"type:varchar(64)"`
	AvatarURL   string    `gorm:"type:varchar(512)"`

	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type UserSecurity struct {
	UserID                uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email                 string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash          string    `gorm:"type:varchar(255);not null"`
	EmailVerified         bool      `gorm:"default:false"`
	RefreshToken          string    `gorm:"type:text"`
	ConfirmationCode      string    `gorm:"type:text"`
	ConfirmationExpiresAt *time.Time
	LastLoginAt           *time.Time
}

func (*User) TableName() string {
	return "users"
}

func (*UserProfile) TableName() string {
	return "user_profile"
}

func (*UserSecurity) TableName() string {
	return "user_security"
}

func ValidateRole(role Role) error {
	if role != UserRole && role != AdminRole {
		return ErrInvalidRole
	}
	return nil
}

func ValidateRoleString(role string) error {
	return ValidateRole(Role(role))
}
