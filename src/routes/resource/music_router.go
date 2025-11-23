package resources

import (
	"magic-server-2026/src/controllers"
	"magic-server-2026/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func MusicRouter(router fiber.Router) {
	api := router.Group("/music")
	api.Get("/", middlewares.ResourceTokenMiddleware, controllers.GetAllMusic)
	api.Get("/:id", middlewares.ResourceTokenMiddleware, controllers.GetMusic)
	api.Post("/vote/:id", middlewares.CSRFTokenMiddleware, controllers.IncrementVote)
	api.Get("/votes/:id", middlewares.ResourceTokenMiddleware, controllers.CanUserVote)
}
