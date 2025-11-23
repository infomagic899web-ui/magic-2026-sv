package resources

import (
	"magic899-server/src/controllers"

	"github.com/gofiber/fiber/v3"
)

func GetPlayerRouter(router fiber.Router) {
	api := router.Group("/player")
	api.Get("/:filename", controllers.GetVideoPlayer)
}
