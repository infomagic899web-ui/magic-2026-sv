package resources

import (
	"magic-server-2026/src/controllers"
	"magic-server-2026/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func YoutubeMagicVideosRouter(router fiber.Router) {
	api := router.Group("/youtube-videos")
	api.Get("/magic/:id", middlewares.ResourceTokenMiddleware, controllers.GetMagicVideosByShowID)
}
