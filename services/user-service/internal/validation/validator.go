package validation

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/domain"
	"regexp"
	"strings"
	"unicode"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("displayname", validateDisplayName)

	return &Validator{validate: v}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validate.Struct(i); err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if len(username) < 3 || len(username) > 64 {
		return false
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_-]*$`, username)
	return matched
}

func validateDisplayName(fl validator.FieldLevel) bool {
	displayName := fl.Field().String()

	if len(displayName) < 1 || len(displayName) > 64 {
		return false
	}

	for _, char := range displayName {
		if unicode.IsControl(char) {
			return false
		}
	}

	return true
}

func (v *Validator) formatValidationError(err error) error {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		var errMsgs []string
		for _, fieldErr := range validationErrors {
			errMsgs = append(errMsgs, formatFieldError(fieldErr))
		}
		return fmt.Errorf("%w: %s", domain.ErrValidationFailed, strings.Join(errMsgs, "; "))
	}

	return err
}

func formatFieldError(fieldErr validator.FieldError) string {
	field := fieldErr.Field()

	switch fieldErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "username":
		return fmt.Sprintf("%s must be 3-64 characters, alphanumeric underscore/hyphen allowed", field)
	case "displayname":
		return fmt.Sprintf("%s must be 1-64 characters", field)
	case "password":
		return fmt.Sprintf("%s must be at least 8 characters and contain uppercase, lowercase, number and special character", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fieldErr.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fieldErr.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "len":
		return fmt.Sprintf("%s must be exactly %d characters", field, fieldErr.Param())
	default:
		return fmt.Sprintf("%s failed validation on '%s'", field, fieldErr.Tag())
	}
}
