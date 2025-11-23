package resources

import (
	"magic899-server/src/controllers"
	"magic899-server/src/handlers"

	"github.com/gofiber/fiber/v3"
)

func NewsRouterResource(app fiber.Router) {
	app.Get("/news", controllers.GetNews, handlers.RotateRSToken)
	app.Get("/news/:id", controllers.GetNewsItem, handlers.RotateRSToken)
}
