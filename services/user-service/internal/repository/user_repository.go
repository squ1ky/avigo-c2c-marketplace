package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/model"
	"gorm.io/gorm"
	"strings"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.Status) error
	UpdateEmailVerified(ctx context.Context, userID uuid.UUID) error
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error

	SaveConfirmationCode(ctx context.Context, userID uuid.UUID, code string, expiresAt time.Time) error
	GetConfirmationCode(ctx context.Context, userID uuid.UUID) (code string, expiresAt time.Time, err error)
	ClearConfirmationCode(ctx context.Context, userID uuid.UUID) error

	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	return r.handleConstraintViolation(err)
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return r.findUser(ctx, "users.id = ?", id)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.findUser(ctx, "Security.email = ?", email)
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return r.findUser(ctx, "users.username = ?", username)
}

func (r *userRepository) GetByEmailOrUsername(ctx context.Context, identifier string) (*model.User, error) {
	var user model.User

	err := r.withRelations(r.db.WithContext(ctx)).
		Where("users.username = ? OR Security.email = ?", identifier, identifier).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateStatus(ctx context.Context, userID uuid.UUID, status model.Status) error {
	result := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("status", status)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) UpdateEmailVerified(ctx context.Context, userID uuid.UUID) error {
	return r.updateSecurity(ctx, userID, map[string]interface{}{
		"email_verified": true,
	})
}

func (r *userRepository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	return r.updateSecurity(ctx, userID, map[string]interface{}{
		"refresh_token": refreshToken,
	})
}

func (r *userRepository) SaveConfirmationCode(ctx context.Context, userID uuid.UUID, code string, expiresAt time.Time) error {
	return r.updateSecurity(ctx, userID, map[string]interface{}{
		"confirmation_code":       code,
		"confirmation_expires_at": expiresAt,
	})
}

func (r *userRepository) GetConfirmationCode(ctx context.Context, userID uuid.UUID) (string, time.Time, error) {
	var security model.UserSecurity

	err := r.db.WithContext(ctx).
		Select("confirmation_code", "confirmation_expires_at").
		Where("user_id = ?", userID).
		First(&security).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", time.Time{}, ErrUserNotFound
		}
		return "", time.Time{}, err
	}

	expiresAt := time.Time{}
	if security.ConfirmationExpiresAt != nil {
		expiresAt = *security.ConfirmationExpiresAt
	}

	return security.ConfirmationCode, expiresAt, nil
}

func (r *userRepository) ClearConfirmationCode(ctx context.Context, userID uuid.UUID) error {
	return r.updateSecurity(ctx, userID, map[string]interface{}{
		"confirmation_code":       nil,
		"confirmation_expires_at": nil,
	})
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	return r.updateSecurity(ctx, userID, map[string]interface{}{
		"last_login_at": time.Now(),
	})
}

// Helpers

func (r *userRepository) withRelations(db *gorm.DB) *gorm.DB {
	return db.Joins("Profile").Joins("Security")
}

func (r *userRepository) findUser(ctx context.Context, condition string, args ...interface{}) (*model.User, error) {
	var user model.User

	err := r.withRelations(r.db.WithContext(ctx)).
		Where(condition, args...).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) updateSecurity(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).
		Model(&model.UserSecurity{}).
		Where("user_id = ?", userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) handleConstraintViolation(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 – unique_violation
		switch {
		case strings.Contains(pgErr.ConstraintName, "username"):
			return ErrUserAlreadyExists
		case strings.Contains(pgErr.ConstraintName, "email"):
			return ErrEmailAlreadyExists
		}
	}

	return err
}
