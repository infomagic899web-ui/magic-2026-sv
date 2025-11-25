package resources

import (
	"magic-server-2026/src/controllers"

	"github.com/gofiber/fiber/v3"
)

func YoutubeMagicVideosRouter(router fiber.Router) {
	api := router.Group("/youtube-videos")
	api.Get("/magic/:id", controllers.GetMagicVideosByShowID)
}
