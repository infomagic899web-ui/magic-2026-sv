package resources

import (
	"magic899-server/src/controllers"
	"magic899-server/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func YoutubeMagicVideosRouter(router fiber.Router) {
	api := router.Group("/youtube-videos")
	api.Get("/magic/:id", middlewares.ResourceTokenMiddleware, controllers.GetMagicVideosByShowID)
}
