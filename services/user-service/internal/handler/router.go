package handler

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler *AuthHandler
	userHandler *UserHandler
}

func NewRouter(authHandler *AuthHandler, userHandler *UserHandler) *Router {
	return &Router{
		authHandler: authHandler,
		userHandler: userHandler,
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

			auth.POST("/change-password", r.authHandler.ChangePassword)
		}

		users := v1.Group("/users")
		{
			users.GET("/me/profile", r.userHandler.GetMyProfile)
			users.PATCH("/me/profile", r.userHandler.UpdateMyProfile)
			users.GET("/:id/profile", r.userHandler.GetUserProfileByID)
		}
	}
}
