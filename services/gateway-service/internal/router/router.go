package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/squ1ky/avigo-c2c-marketplace/services/gateway-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/gateway-service/internal/middleware"
	"github.com/squ1ky/avigo-c2c-marketplace/services/gateway-service/internal/proxy"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}))

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

		listingProxy := proxy.ReverseProxy(cfg.Services.ListingServiceURL)
		orderProxy := listingProxy
		publicListings := v1.Group("/listings")
		{
			publicListings.GET("/:id", listingProxy)
		}
		v1.GET("/categories", listingProxy)

		// JWT require
		protected := v1.Group("")
		protected.Use(middleware.JWTMiddleware(cfg.JWT.Secret))
		{
			userProxy := proxy.ReverseProxy(cfg.Services.UserServiceURL)

			authProtected := protected.Group("/auth")
			{
				authProtected.POST("/logout", userProxy)
				authProtected.POST("/refresh", userProxy)
				authProtected.GET("/me", userProxy)
				authProtected.POST("/change-password", userProxy)
			}

			usersProtected := protected.Group("/users")
			{
				usersProtected.GET("/me/profile", userProxy)
				usersProtected.PATCH("/me/profile", userProxy)
			}

			listingsProtected := protected.Group("/listings")
			{
				listingsProtected.GET("/my", listingProxy)
				listingsProtected.POST("", listingProxy)
				listingsProtected.PUT("/:id", listingProxy)
				listingsProtected.DELETE("/:id", listingProxy)
			}

			ordersProtected := protected.Group("/orders")
			{
				ordersProtected.GET("/purchases", orderProxy)
				ordersProtected.GET("/sales", orderProxy)
			}

			mediaProtected := protected.Group("/media")
			{
				mediaProtected.POST("/upload", listingProxy)
			}
		}
	}

	return router
}
