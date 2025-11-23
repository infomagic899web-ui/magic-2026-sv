package middlewares

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App) {
	app.Use(CORSMiddleware)
	app.Use(CookieStealBlocker())

	env := os.Getenv("ENV")
	if env == "production" {
		app.Use(SilentErrorsMiddleware)
	}

	fmt.Println("Middleware Initialized")
}
