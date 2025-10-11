package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/domain"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	config.JWTConfig
}

func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		JWTConfig: *cfg,
	}
}

func (m *JWTManager) GenerateAccessToken(userID, email, role string) (string, error) {
	return m.generateToken(userID, email, role, m.AccessTokenDuration)
}

func (m *JWTManager) GenerateRefreshToken(userID, email, role string) (string, error) {
	return m.generateToken(userID, email, role, m.RefreshTokenDuration)
}

func (m *JWTManager) generateToken(userID, email, role string, duration time.Duration) (string, error) {
	if err := domain.ValidateRoleString(role); err != nil {
		return "", err
	}

	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(m.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.Secret), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if err := domain.ValidateRoleString(claims.Role); err != nil {
		return nil, err
	}

	return claims, nil
}
