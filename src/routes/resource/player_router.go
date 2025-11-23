package resources

import (
	"magic-server-2026/src/controllers"

	"github.com/gofiber/fiber/v3"
)

func GetPlayerRouter(router fiber.Router) {
	api := router.Group("/player")
	api.Get("/:filename", controllers.GetVideoPlayer)
}
