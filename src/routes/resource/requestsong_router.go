package resources

import (
	"magic899-server/src/controllers"
	"magic899-server/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func RequestedSongRouter(router fiber.Router) {
	api := router.Group("/requestsongs")
	api.Get("/", controllers.GetAllRequestSongs)
	api.Get("/:id", controllers.GetRequestSong)
	api.Post("/", middlewares.RateLimiterMiddleware(), middlewares.CSRFTokenMiddleware, controllers.CreateRequestSong)
	api.Put("/name/:name", middlewares.RateLimiterMiddleware(), middlewares.CSRFTokenMiddleware, controllers.UpdateRequestSong)
	api.Get("/request/eligibility", middlewares.RateLimiterMiddleware(), middlewares.CSRFTokenMiddleware, controllers.CheckEligibility)
}
