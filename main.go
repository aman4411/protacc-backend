package main

import (
	"github.com/aman4411/protacc-backend/internal/db"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	if err := godotenv.Load(); err != nil {
		log.Println("ERROR: No .env file found, using environment variables")
	}
	log.Println("INFO: Current Profile =", os.Getenv("ENVIRONMENT"))

	// Initialize database connection
	db.InitDB()
	defer db.CloseDB()

	// Set up Fiber app
	app := fiber.New()
	app.Use(cors.New())

	// Test route
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Graceful shutdown on SIGINT/SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		db.CloseDB()
		os.Exit(0)
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start listening
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
