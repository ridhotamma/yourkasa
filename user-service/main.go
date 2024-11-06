package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/user-service/config"
	"github.com/ridhotamma/yourkasa/user-service/routes"
)

func main() {
	db := config.InitDB()
	r := gin.Default()

	routes.SetupRoutes(r, db)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
