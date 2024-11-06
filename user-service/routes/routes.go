package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/user-service/controllers"
	"github.com/ridhotamma/yourkasa/user-service/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	userController := controllers.NewUserController(db)

	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", userController.GetCurrentUser)

			admin := users.Group("/")
			admin.Use(middleware.RequireRole("admin"))
			{
				users.GET("/:id", userController.GetByID)
				admin.POST("/", userController.Create)
				admin.GET("/", userController.List)
				admin.PUT("/:id", userController.Update)
				admin.DELETE("/:id", userController.Delete)
			}
		}
	}
}
