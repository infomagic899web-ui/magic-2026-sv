package resources

import (
	"magic-server-2026/src/controllers"
	"magic-server-2026/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func ShowRouter(router fiber.Router) {
	api := router.Group("/shows")
	api.Get("/", middlewares.ResourceTokenMiddleware, controllers.GetShows)
	api.Get("/:id", middlewares.ResourceTokenMiddleware, controllers.GetShow)
	api.Get("/name/:showName", middlewares.ResourceTokenMiddleware, controllers.GetByShowName)

}
