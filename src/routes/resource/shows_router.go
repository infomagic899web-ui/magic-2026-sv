package resources

import (
	"magic899-server/src/controllers"
	"magic899-server/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func ShowRouter(router fiber.Router) {
	api := router.Group("/shows")
	api.Get("/", middlewares.ResourceTokenMiddleware, controllers.GetShows)
	api.Get("/:id", middlewares.ResourceTokenMiddleware, controllers.GetShow)
	api.Get("/name/:showName", middlewares.ResourceTokenMiddleware, controllers.GetByShowName)

}
