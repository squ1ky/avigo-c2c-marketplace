package http

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()

	handler := NewHandler()

	router.GET("/health", handler.HealthCheck)

	return router
}
