package resources

import (
	"magic899-server/src/controllers"
	"magic899-server/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func MusicRouter(router fiber.Router) {
	api := router.Group("/music")
	api.Get("/", middlewares.ResourceTokenMiddleware, controllers.GetAllMusic)
	api.Get("/:id", middlewares.ResourceTokenMiddleware, controllers.GetMusic)
	api.Post("/vote/:id", middlewares.CSRFTokenMiddleware, controllers.IncrementVote)
	api.Get("/votes/:id", middlewares.ResourceTokenMiddleware, controllers.CanUserVote)
}
