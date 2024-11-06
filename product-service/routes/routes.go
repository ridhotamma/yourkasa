package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/product-service/controllers"
	"github.com/ridhotamma/yourkasa/product-service/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	productController := controllers.NewProductController(db)

	api := r.Group("/api/v1")
	{
		products := api.Group("/products")
		products.Use(middleware.AuthMiddleware())
		{
			// Public routes (require authentication)
			products.GET("/:id", productController.GetByID)
			products.GET("/", productController.List)

			// Admin/Owner only routes
			authorized := products.Group("/")
			authorized.Use(middleware.RequireRole("admin", "owner"))
			{
				authorized.POST("/", productController.Create)
				authorized.PUT("/:id", productController.Update)
				authorized.DELETE("/:id", productController.Delete)
			}
		}
	}
}
