package model

import (
	"time"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type UserSecurity struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id" json:"user_id"`
	LastLoginAt time.Time

	PasswordHash string `gorm:"type:text;not_null"`

	User User `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (UserSecurity) TableName() string {
	return "user_security"
}
