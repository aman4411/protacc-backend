package main

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

func main() {
	app := fiber.New()

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := app.Listen(":" + port)
	if err != nil {
		return
	}
}
