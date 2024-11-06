package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/auth-service/controllers"
	"github.com/ridhotamma/yourkasa/auth-service/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)

			// Protected routes
			protected := auth.Group("/")
			protected.Use(middleware.AuthMiddleware())
			{
				protected.POST("/change-password", authController.ChangePassword)
			}
		}
	}
}
