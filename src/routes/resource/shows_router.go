package resources

import (
	"magic-server-2026/src/controllers"

	"github.com/gofiber/fiber/v3"
)

func ShowRouter(router fiber.Router) {
	api := router.Group("/shows")
	api.Get("/", controllers.GetShows)
	api.Get("/:id", controllers.GetShow)
	api.Get("/name/:showName", controllers.GetByShowName)

}
