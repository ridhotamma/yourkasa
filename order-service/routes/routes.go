package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/order-service/controllers"
	"github.com/ridhotamma/yourkasa/order-service/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize controllers
	checkoutController := controllers.NewCheckoutController(db)
	orderController := controllers.NewOrderController(db)

	api := r.Group("/api/v1")
	{
		// Checkout/Cart routes
		cart := api.Group("/cart")
		cart.Use(middleware.AuthMiddleware())
		{
			cart.POST("/items", checkoutController.AddToCart)
			cart.GET("/items", checkoutController.GetCart)
			cart.PUT("/items/:id", checkoutController.UpdateCartItem)
			cart.DELETE("/items/:id", checkoutController.RemoveFromCart)
		}

		// Order routes
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("/", orderController.Create)
			orders.GET("/", orderController.List)
			orders.GET("/:id", orderController.GetByID)
			orders.POST("/:id/cancel", orderController.Cancel)
		}
	}
}
