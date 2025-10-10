package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/auth"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/kafka"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/model"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/repository"
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
	userRepo             repository.UserRepository
	txManager            repository.TransactionManager
	jwtManager           *auth.JWTManager
	kafkaProducer        *kafka.Producer
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(
	userRepo repository.UserRepository,
	txManager repository.TransactionManager,
	jwtManager *auth.JWTManager,
	kafkaProducer *kafka.Producer,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:             userRepo,
		txManager:            txManager,
		jwtManager:           jwtManager,
		kafkaProducer:        kafkaProducer,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

type RegisterRequest struct {
	Username    string
	Email       string
	Password    string
	DisplayName string
	Phone       string
}

type RegisterResponse struct {
	UserID string
	Email  string
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	confirmationCode, err := generateConfirmationCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate confirmation code: %w", err)
	}

	expiresAt := time.Now().Add(15 * time.Minute)

	user := model.User{
		ID:       uuid.New(),
		Username: req.Username,
		Status:   model.StatusPending,
		Role:     model.UserRole,
		Profile: &model.UserProfile{
			DisplayName: req.DisplayName,
			Phone:       req.Phone,
		},
		Security: &model.UserSecurity{
			Email:                 req.Email,
			PasswordHash:          string(passwordHash),
			EmailVerified:         false,
			ConfirmationCode:      confirmationCode,
			ConfirmationExpiresAt: &expiresAt,
		},
	}

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

	s.kafkaProducer.SendEmailVerification(event)

	return &RegisterResponse{
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

type LoginRequest struct {
	Login    string // Username/Email
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	User         *UserInfo
}

type UserInfo struct {
	ID          string
	Username    string
	Email       string
	DisplayName string
	Role        string
	Status      string
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmailOrUsername(ctx, req.Login)
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

	s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &UserInfo{
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
	savedCode, expiresAt, err := s.userRepo.GetConfirmationCode(ctx, userID)
	if err != nil {
		return err
	}

	if !expiresAt.IsZero() && time.Now().After(expiresAt) {
		return ErrConfirmationExpired
	}

	if savedCode != code {
		return ErrInvalidConfirmation
	}

	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tx := repository.GetTxFromContext(txCtx)

		txRepo := s.userRepo.WithTx(tx)

		if err := txRepo.UpdateEmailVerified(ctx, userID); err != nil {
			return err
		}

		if err := txRepo.UpdateStatus(ctx, userID, model.StatusActive); err != nil {
			return err
		}

		if err := txRepo.ClearConfirmationCode(txCtx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*LoginResponse, error) {
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

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User: &UserInfo{
			ID:          user.ID.String(),
			Username:    user.Username,
			Email:       user.Security.Email,
			DisplayName: user.Profile.DisplayName,
			Role:        string(user.Role),
			Status:      string(user.Status),
		},
	}, nil
}
