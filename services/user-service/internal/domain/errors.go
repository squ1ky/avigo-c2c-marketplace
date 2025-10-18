package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user with this username already exists")
	ErrEmailAlreadyExists = errors.New("user with this email already exists")

	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailNotVerified    = errors.New("email not verified")
	ErrInvalidConfirmation = errors.New("invalid confirmation code")
	ErrConfirmationExpired = errors.New("confirmation code expired")
	ErrUserAlreadyVerified = errors.New("user already verified")

	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")

	ErrInvalidRole      = errors.New("invalid role")
	ErrValidationFailed = errors.New("validation failed")
)
