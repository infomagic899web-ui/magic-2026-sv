package resources

import (
	"magic-server-2026/src/controllers"
	"magic-server-2026/src/middlewares"

	"github.com/gofiber/fiber/v3"
)

func RequestedSongRouter(router fiber.Router) {
	api := router.Group("/requestsongs")
	api.Get("/", middlewares.ResourceTokenMiddleware, controllers.GetAllRequestSongs)
	api.Get("/:id", middlewares.ResourceTokenMiddleware, controllers.GetRequestSong)
	api.Post("/", middlewares.RateLimiterMiddleware(), middlewares.CSRFTokenMiddleware, controllers.CreateRequestSong)
	api.Put("/name/:name", middlewares.RateLimiterMiddleware(), middlewares.CSRFTokenMiddleware, controllers.UpdateRequestSong)
	api.Get("/request/eligibility", middlewares.RateLimiterMiddleware(), middlewares.ResourceTokenMiddleware, controllers.CheckEligibility)
}
