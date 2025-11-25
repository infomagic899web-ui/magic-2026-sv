package resources

import (
	"magic-server-2026/src/controllers"

	"github.com/gofiber/fiber/v3"
)

func MoviesRouter(router fiber.Router) {
	api := router.Group("/movies")
	api.Get("/", controllers.GetMovies)
	api.Get("/:id", controllers.GetMovie)
	api.Get("/latest/now", controllers.GetLatestMovie)
}
