package server

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
)

func Start(app *fiber.App) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := app.Listen(":" + port); err != nil {
		panic(fmt.Sprintf("Error starting server: %v", err))
	}
}
