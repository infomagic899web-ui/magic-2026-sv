package resources

import (
	"magic-server-2026/src/controllers"
	"magic-server-2026/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func MagicVideoRouter(router fiber.Router) {
	api := router.Group("/magic-videos")
	api.Get("/", middlewares.ResourceTokenMiddleware, controllers.GetMagicVideos)
	api.Get("/:id", middlewares.ResourceTokenMiddleware, controllers.GetMagicVideo)
}
