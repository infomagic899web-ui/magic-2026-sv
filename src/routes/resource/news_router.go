package resources

import (
	"magic-server-2026/src/controllers"
	"magic-server-2026/src/handlers"

	"github.com/gofiber/fiber/v3"
)

func NewsRouterResource(app fiber.Router) {
	app.Get("/news", controllers.GetNews, handlers.RotateRSToken)
	app.Get("/news/:id", controllers.GetNewsItem, handlers.RotateRSToken)
}
