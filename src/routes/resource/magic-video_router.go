package resources

import (
	"magic-server-2026/src/controllers"

	"github.com/gofiber/fiber/v3"
)

func MagicVideoRouter(router fiber.Router) {
	api := router.Group("/magic-videos")
	api.Get("/", controllers.GetMagicVideos)
	api.Get("/:id", controllers.GetMagicVideo)
}
