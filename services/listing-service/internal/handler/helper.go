package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) getTargetUserID(c *gin.Context) (uuid.UUID, error) {
	idStr := c.Param("user_id")
	if idStr == "" {
		return uuid.Nil, errors.New("empy user id")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid uuid format")
	}

	return id, nil
}

func getUserIDFromHeader(c *gin.Context) (uuid.UUID, bool) {
	idStr := c.GetHeader("X-User-ID")
	if idStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "missing authentication",
			"code":  "UNAUTHORIZED",
		})
		return uuid.Nil, false
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
			"code":  "INVALID_TOKEN_PAYLOAD",
		})
		return uuid.Nil, false
	}

	return id, true
}

func invalidJSON(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error":   "invalid request body",
		"code":    "INVALID_JSON",
		"details": err.Error(),
	})
}

func badRequest(c *gin.Context, code string, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
		"code":  code,
	})
}
