package routes

import (
	"magic-server-2026/src/middlewares"
	resources "magic-server-2026/src/routes/resource"
	tokens "magic-server-2026/src/routes/token"
	"time"

	"github.com/gofiber/fiber/v3"
)

func SetupRouter(app *fiber.App) {

	rl := middlewares.NewRateLimiter(12 * time.Hour)

	api := app.Group("/api", middlewares.ReferrerMiddleware, middlewares.MakeTrustedRateLimiterMiddleware(rl))

	apiSecure := api.Group("/v1", middlewares.ResourceTokenMiddleware)

	resourceRoutes := []func(router fiber.Router){
		resources.NewsRouterResource,
		resources.MusicRouter,
		resources.ShowRouter,
		resources.YoutubeMagicVideosRouter,
		resources.MagicVideoRouter,
		resources.MoviesRouter,
		resources.RequestedSongRouter,
		resources.ShoutboxMailerRouter,
	}

	for _, r := range resourceRoutes {
		r(apiSecure)
	}

	resources.TestRouter(api)
	resources.ImageRouter(api)
	resources.GetPlayerRouter(api)
	tokens.SetupTokenRouter(api)
}
