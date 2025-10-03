package model

import (
	"time"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type UserProfile struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id" json:"user_id"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	DisplayName string `gorm:"type:varchar(64);not null" json:"display_name"`
	About       string `gorm:"type:text" json:"about"`
	Country     string `gorm:"type:varchar(64);not null" json:"country"`
	City        string `gorm:"type:varchar(64);not null" json:"city"`
	AvatarURL   string `gorm:"type:text;column:avatar_url" json:"avatar_url"`
	Email       string `gorm:"type:varchar(255)" json:"email"`
	Phone       string `gorm:"type:varchar(32)" json:"phone"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
