package resources

import (
	"magic-server-2026/src/controllers"
	"magic-server-2026/src/handlers"
	"magic-server-2026/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func NewsRouterResource(app fiber.Router) {
	app.Get("/news", middlewares.ResourceTokenMiddleware, controllers.GetNews, handlers.RotateRSToken)
	app.Get("/news/:id", middlewares.ResourceTokenMiddleware, controllers.GetNewsItem, handlers.RotateRSToken)
}
