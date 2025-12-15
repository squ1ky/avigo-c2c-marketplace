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
		userProxy := proxy.ReverseProxy(cfg.Services.UserServiceURL)
		listingProxy := proxy.ReverseProxy(cfg.Services.ListingServiceURL)
		orderProxy := listingProxy

		auth := v1.Group("/auth")
		{
			auth.POST("/register", userProxy)
			auth.POST("/login", userProxy)
			auth.POST("/confirm-email", userProxy)
		}

		v1.GET("/categories", listingProxy)
		v1.GET("/account/:user_id/listings", listingProxy)
		v1.GET("/account/:user_id/listings/:listing_id", listingProxy)

		protected := v1.Group("")
		protected.Use(middleware.JWTMiddleware(cfg.JWT.Secret))
		{
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
				usersProtected.GET("/:id/profile", userProxy)
			}

			accountProtected := protected.Group("/account/:user_id")
			{
				accountProtected.POST("/listings", listingProxy)
				accountProtected.PUT("/listings/:listing_id", listingProxy)
				accountProtected.DELETE("/listings/:listing_id", listingProxy)

				accountProtected.GET("/orders/purchases", orderProxy)
				accountProtected.GET("/orders/sales", orderProxy)
			}

			mediaProtected := protected.Group("/media")
			{
				mediaProtected.POST("/upload", listingProxy)
			}
		}
	}

	return router
}
