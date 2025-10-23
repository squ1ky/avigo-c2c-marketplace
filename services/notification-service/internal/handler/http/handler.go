package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	// пока оставил пустым, так как сейчас все необходимое обрабатывает kafka handler
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "notification-service",
	})
}
