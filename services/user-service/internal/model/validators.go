package model

import "errors"

var ErrInvalidRole = errors.New("invalid role")

func ValidateRole(role Role) error {
	if role != UserRole && role != AdminRole {
		return ErrInvalidRole
	}
	return nil
}

func ValidateRoleString(role string) error {
	return ValidateRole(Role(role))
}
