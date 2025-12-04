package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

const headerUserID = "X-User-ID"

func getUserIDFromHeader(c *gin.Context) (uuid.UUID, bool) {
	userIDHeader := c.GetHeader(headerUserID)
	if userIDHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDHeader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID format",
			"code":  "INVALID_USER_ID",
		})
		return uuid.Nil, false
	}

	return userID, true
}
