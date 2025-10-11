package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/auth"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/domain"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/dto"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/kafka"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/repository"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/validation"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrEmailNotVerified    = errors.New("email not verified")
	ErrInvalidConfirmation = errors.New("invalid confirmation code")
	ErrConfirmationExpired = errors.New("confirmation code expired")
	ErrUserAlreadyVerified = errors.New("user already verified")
)

type AuthService struct {
	userRepo               repository.UserRepository
	txManager              repository.TransactionManager
	validator              *validation.Validator
	jwtManager             *auth.JWTManager
	kafkaProducer          *kafka.Producer
	confirmationCodeExpiry time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	txManager repository.TransactionManager,
	validator *validation.Validator,
	jwtManager *auth.JWTManager,
	kafkaProducer *kafka.Producer,
	confirmationCodeExpiry time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:               userRepo,
		txManager:              txManager,
		validator:              validator,
		jwtManager:             jwtManager,
		kafkaProducer:          kafkaProducer,
		confirmationCodeExpiry: confirmationCodeExpiry,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	confirmationCode, err := generateConfirmationCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate confirmation code: %w", err)
	}

	expiresAt := time.Now().UTC().Add(s.confirmationCodeExpiry)

	user := domain.User{
		ID:       uuid.New(),
		Username: req.Username,
		Status:   domain.StatusPending,
		Role:     domain.UserRole,
		Profile: &domain.UserProfile{
			DisplayName: req.DisplayName,
			Phone:       req.Phone,
			Country:     req.Country,
			City:        req.City,
		},
		Security: &domain.UserSecurity{
			Email:                 req.Email,
			PasswordHash:          string(passwordHash),
			EmailVerified:         false,
			ConfirmationCode:      confirmationCode,
			ConfirmationExpiresAt: &expiresAt,
		},
	}

	// We will receive err here if user already exists
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}

	event := kafka.NewEmailVerificationEvent(
		user.ID.String(),
		user.Security.Email,
		user.Profile.DisplayName,
		confirmationCode,
		expiresAt,
	)

	// TODO: handle err when notification-service can't
	s.kafkaProducer.SendEmailVerification(event)

	return &dto.RegisterResponse{
		UserID: user.ID.String(),
		Email:  user.Security.Email,
	}, nil
}

func generateConfirmationCode() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)

	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmailOrUsername(ctx, req.Identifier)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Security.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.Security.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(
		user.ID.String(),
		user.Security.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(
		user.ID.String(),
		user.Security.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, err
	}

	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserInfo{
			ID:          user.ID.String(),
			Username:    user.Username,
			Email:       user.Security.Email,
			DisplayName: user.Profile.DisplayName,
			Role:        string(user.Role),
			Status:      string(user.Status),
		},
	}, nil
}

func (s *AuthService) ConfirmEmail(ctx context.Context, userID uuid.UUID, code string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Security.EmailVerified {
		return ErrUserAlreadyVerified
	}

	savedCode, err := s.userRepo.GetConfirmationCode(ctx, userID)
	if err != nil {
		return err
	}

	if !savedCode.ExpiresAt.IsZero() && time.Now().After(savedCode.ExpiresAt) {
		return ErrConfirmationExpired
	}

	if savedCode.Code != code {
		return ErrInvalidConfirmation
	}

	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tx := repository.GetTxFromContext(txCtx)

		txRepo := s.userRepo.WithTx(tx)

		if err := txRepo.UpdateEmailVerified(txCtx, userID); err != nil {
			return err
		}

		if err := txRepo.UpdateStatus(txCtx, userID, domain.StatusActive); err != nil {
			return err
		}

		if err := txRepo.ClearConfirmationCode(txCtx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Security.RefreshToken == "" || user.Security.RefreshToken != refreshToken {
		return nil, auth.ErrInvalidToken
	}

	newAccessToken, err := s.jwtManager.GenerateAccessToken(
		user.ID.String(),
		user.Security.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(
		user.ID.String(),
		user.Security.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.UpdateRefreshToken(ctx, userID, newRefreshToken); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User: dto.UserInfo{
			ID:          user.ID.String(),
			Username:    user.Username,
			Email:       user.Security.Email,
			DisplayName: user.Profile.DisplayName,
			Role:        string(user.Role),
			Status:      string(user.Status),
		},
	}, nil
}
