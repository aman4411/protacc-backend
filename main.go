package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Important: check for error
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
