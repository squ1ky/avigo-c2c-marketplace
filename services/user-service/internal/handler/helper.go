package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InvalidJSON(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error":   "invalid request body",
		"code":    "INVALID_JSON",
		"details": err.Error(),
	})
}
