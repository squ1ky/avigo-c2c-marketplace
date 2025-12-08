package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/squ1ky/avigo-c2c-marketplace/services/listing-service/internal/service"
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
		statusCode, errorCode, msg := mapDomainError(err)

		c.JSON(statusCode, ErrorResponse{
			Error: msg,
			Code:  errorCode,
		})
	}
}

func mapDomainError(err error) (statusCode int, code string, message string) {
	switch {
	case errors.Is(err, service.ErrListingNotFound):
		return http.StatusNotFound, "LISTING_NOT_FOUND", err.Error()
	case errors.Is(err, service.ErrListingNotActive):
		return http.StatusBadRequest, "LISTING_NOT_ACTIVE", err.Error()
	case errors.Is(err, service.ErrListingSold):
		return http.StatusConflict, "LISTING_ALREADY_SOLD", err.Error()
	case errors.Is(err, service.ErrSelfPurchase):
		return http.StatusBadRequest, "SELF_PURCHASE_FORBIDDEN", err.Error()

	case errors.Is(err, service.ErrOrderNotFound):
		return http.StatusNotFound, "ORDER_NOT_FOUND", err.Error()
	case errors.Is(err, service.ErrInvalidStatus):
		return http.StatusBadRequest, "INVALID_ORDER_STATUS", err.Error()
	case errors.Is(err, service.ErrForbidden):
		return http.StatusForbidden, "FORBIDDEN", err.Error()

	case errors.Is(err, service.ErrOrderNotCompleted):
		return http.StatusBadRequest, "ORDER_NOT_COMPLETED", err.Error()
	case errors.Is(err, service.ErrReviewAlreadyExists):
		return http.StatusConflict, "REVIEW_ALREADY_EXISTS", err.Error()
	case errors.Is(err, service.ErrReviewForbidden):
		return http.StatusForbidden, "REVIEW_FORBIDDEN", err.Error()

	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"
	}
}
