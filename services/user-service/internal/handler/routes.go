package handler

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler *AuthHandler
}

func NewRouter(authHandler *AuthHandler) *Router {
	return &Router{
		authHandler: authHandler,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/confirm-email", r.authHandler.ConfirmEmail)

			// gateway already checked JWT
			auth.POST("/logout", r.authHandler.Logout)
			auth.POST("/refresh", r.authHandler.RefreshToken)

			auth.GET("/me", r.authHandler.Me)
		}
	}
}
