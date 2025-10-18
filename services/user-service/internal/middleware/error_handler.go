package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/domain"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		statusCode, errorCode, message := mapDomainError(err)

		if c.Writer.Status() != http.StatusOK {

		}

		c.JSON(statusCode, ErrorResponse{
			Error: message,
			Code:  errorCode,
		})
	}
}

func mapDomainError(err error) (statusCode int, code string, message string) {
	switch {
	// User errors
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, "USER_NOT_FOUND", err.Error()

	case errors.Is(err, domain.ErrUserAlreadyExists) || errors.Is(err, domain.ErrEmailAlreadyExists):
		return http.StatusConflict, "USER_ALREADY_EXISTS", err.Error()

	// Auth errors
	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password"

	case errors.Is(err, domain.ErrEmailNotVerified):
		return http.StatusForbidden, "EMAIL_NOT_VERIFIED", err.Error()

	case errors.Is(err, domain.ErrInvalidConfirmation):
		return http.StatusBadRequest, "INVALID_CONFIRMATION", err.Error()

	case errors.Is(err, domain.ErrConfirmationExpired):
		return http.StatusBadRequest, "CONFIRMATION_EXPIRED", err.Error()

	case errors.Is(err, domain.ErrUserAlreadyVerified):
		return http.StatusConflict, "USER_ALREADY_VERIFIED", err.Error()

	// Token errors
	case errors.Is(err, domain.ErrInvalidToken):
		return http.StatusUnauthorized, "INVALID_TOKEN", err.Error()

	case errors.Is(err, domain.ErrExpiredToken):
		return http.StatusUnauthorized, "EXPIRED_TOKEN", err.Error()

	// Validation errors
	case errors.Is(err, domain.ErrValidationFailed):
		return http.StatusBadRequest, "VALIDATION_FAILED", err.Error()

	case errors.Is(err, domain.ErrInvalidRole):
		return http.StatusBadRequest, "INVALID_ROLE", err.Error()

	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"
	}
}
