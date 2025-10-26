package domain

import "errors"

var (
	ErrMissingUserID           = errors.New("user_id is required")
	ErrMissingEmail            = errors.New("email is required")
	ErrMissingConfirmationCode = errors.New("confirmation_code is required")
	ErrExpiredCode             = errors.New("confirmation code has expired")
	ErrInvalidUUID             = errors.New("invalid UUID format")
)
