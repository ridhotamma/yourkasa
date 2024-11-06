package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/product-service/controllers"
	"github.com/ridhotamma/yourkasa/product-service/middleware"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize all controllers
	productController := controllers.NewProductController(db)
	categoryController := controllers.NewCategoryController(db)
	groupController := controllers.NewGroupController(db)
	variantController := controllers.NewVariantController(db)
	addonController := controllers.NewAddonController(db)

	api := r.Group("/api/v1")
	{
		// Product routes
		products := api.Group("/products")
		products.Use(middleware.AuthMiddleware())
		{
			// Public routes (require authentication)
			products.GET("/:id", productController.GetByID)
			products.GET("/", productController.List)

			// Admin/Owner only routes
			authorizedProducts := products.Group("/")
			authorizedProducts.Use(middleware.RequireRole("admin", "owner"))
			{
				authorizedProducts.POST("/", productController.Create)
				authorizedProducts.PUT("/:id", productController.Update)
				authorizedProducts.DELETE("/:id", productController.Delete)
			}
		}

		// Category routes
		categories := api.Group("/categories")
		categories.Use(middleware.AuthMiddleware())
		{
			// Public routes
			categories.GET("/:id", categoryController.GetByID)
			categories.GET("/", categoryController.List)

			// Admin/Owner only routes
			authorizedCategories := categories.Group("/")
			authorizedCategories.Use(middleware.RequireRole("admin", "owner"))
			{
				authorizedCategories.POST("/", categoryController.Create)
				authorizedCategories.PUT("/:id", categoryController.Update)
				authorizedCategories.DELETE("/:id", categoryController.Delete)
			}
		}

		// Group routes
		groups := api.Group("/groups")
		groups.Use(middleware.AuthMiddleware())
		{
			// Public routes
			groups.GET("/:id", groupController.GetByID)
			groups.GET("/", groupController.List)

			// Admin/Owner only routes
			authorizedGroups := groups.Group("/")
			authorizedGroups.Use(middleware.RequireRole("admin", "owner"))
			{
				authorizedGroups.POST("/", groupController.Create)
				authorizedGroups.PUT("/:id", groupController.Update)
				authorizedGroups.DELETE("/:id", groupController.Delete)
				authorizedGroups.POST("/:id/products", groupController.AddProductToGroup)
				authorizedGroups.DELETE("/:id/products/:productId", groupController.RemoveProductFromGroup)
			}
		}

		// Variant routes
		variants := api.Group("/variants")
		variants.Use(middleware.AuthMiddleware())
		{
			// Public routes
			variants.GET("/:id", variantController.GetByID)
			variants.GET("/product/:productId", variantController.GetByProductID)

			// Admin/Owner only routes
			authorizedVariants := variants.Group("/")
			authorizedVariants.Use(middleware.RequireRole("admin", "owner"))
			{
				authorizedVariants.POST("/", variantController.Create)
				authorizedVariants.PUT("/:id", variantController.Update)
				authorizedVariants.DELETE("/:id", variantController.Delete)
			}
		}

		// Addon routes
		addons := api.Group("/addons")
		addons.Use(middleware.AuthMiddleware())
		{
			// Public routes
			addons.GET("/:id", addonController.GetByID)
			addons.GET("/product/:productId", addonController.GetByProductID)

			// Admin/Owner only routes
			authorizedAddons := addons.Group("/")
			authorizedAddons.Use(middleware.RequireRole("admin", "owner"))
			{
				authorizedAddons.POST("/", addonController.Create)
				authorizedAddons.PUT("/:id", addonController.Update)
				authorizedAddons.DELETE("/:id", addonController.Delete)
			}
		}
	}
}
