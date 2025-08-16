package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"task-manager-backend/config"
	"task-manager-backend/routes"
	"task-manager-backend/utils"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Could not load config: %v", err)
	}

	// Initialize the database connection
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Could not connect to the database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	// Set up Gin router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	// Set up routes
	routes.SetupRoutes(r, db)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("❌ Could not start server: %v", err)
	}
}