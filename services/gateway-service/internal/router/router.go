package router

import (
	"github.com/gin-gonic/gin"
	"github.com/squ1ky/avigo-c2c-marketplace/services/gateway-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/gateway-service/internal/middleware"
	"github.com/squ1ky/avigo-c2c-marketplace/services/gateway-service/internal/proxy"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			userProxy := proxy.ReverseProxy(cfg.Services.UserServiceURL)
			auth.POST("/register", userProxy)
			auth.POST("/login", userProxy)
			auth.POST("/confirm-email", userProxy)
		}

		// JWT require
		protected := v1.Group("")
		protected.Use(middleware.JWTMiddleware(cfg.JWT.Secret))
		{
			authProtected := protected.Group("/auth")
			{
				userProxy := proxy.ReverseProxy(cfg.Services.UserServiceURL)
				authProtected.POST("/logout", userProxy)
				authProtected.POST("/refresh", userProxy)
				authProtected.GET("/me", userProxy)
			}
		}
	}

	return router
}
