package main

import (
	"log"
	"os"

	"github.com/aman4411/protacc-backend/app"
	"github.com/aman4411/protacc-backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Println("ERROR: No .env file found, using environment variables")
	}
	log.Println("INFO: Current Profile =", os.Getenv("ENVIRONMENT"))

	// Initialize all dependencies
	factories, err := app.NewFactories()
	if err != nil {
		log.Fatalf("Failed to initialize factories: %v", err)
	}

	// Ensure DB connection is closed on exit
	defer db.CloseDB()

	// Set up Fiber app
	fiberApp := fiber.New()

	// Setup routes
	app.SetupRoutes(fiberApp, factories)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := fiberApp.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
