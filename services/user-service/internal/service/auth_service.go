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

	// TODO: handle err
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
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Security.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.Security.EmailVerified {
		return nil, domain.ErrEmailNotVerified
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
		return domain.ErrUserAlreadyVerified
	}

	savedCode, err := s.userRepo.GetConfirmationCode(ctx, userID)
	if err != nil {
		return err
	}

	if !savedCode.ExpiresAt.IsZero() && time.Now().After(savedCode.ExpiresAt) {
		return domain.ErrConfirmationExpired
	}

	if savedCode.Code != code {
		return domain.ErrInvalidConfirmation
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
		return nil, domain.ErrInvalidToken
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

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Profile == nil || user.Security == nil {
		return nil, domain.ErrUserNotFound
	}

	return &dto.UserInfo{
		ID:          user.ID.String(),
		Username:    user.Username,
		Email:       user.Security.Email,
		DisplayName: user.Profile.DisplayName,
		Role:        string(user.Role),
		Status:      string(user.Status),
	}, nil
}

func (s *AuthService) ChangePassword(
	ctx context.Context,
	userID uuid.UUID,
	req dto.ChangePasswordRequest,
) error {
	if err := s.validator.Validate(req); err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Security.PasswordHash),
		[]byte(req.CurrentPassword),
	); err != nil {
		return domain.ErrInvalidCredentials
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.UpdatePasswordHash(ctx, userID, string(newHash)); err != nil {
		return err
	}

	// invalidate refresh-token for logout on other devices
	_ = s.userRepo.UpdateRefreshToken(ctx, userID, "")

	return nil
}
